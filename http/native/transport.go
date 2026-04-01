package native

import (
	"bytes"
	"context"
	"errors"
	"github.com/mono83/transport"
	"github.com/mono83/transport/http/log"
	"io"
	http2 "net/http"
	"sync"
	"time"
)

// nativeTransport implements [http.Transport] using [net/http].
type nativeTransport struct {
	transport *http2.Transport
	sink      log.Sink
	options   []any
}

func (n nativeTransport) ExecuteRequest(
	ctx context.Context,
	method string,
	url string,
	reqData io.ReadCloser,
	options ...any,
) (int, map[string][]string, io.ReadCloser, error) {
	// Preparing data for logging
	var reqBuf bytes.Buffer
	if n.sink != nil && reqData != nil {
		reqData = teeReadCloser(reqData, &reqBuf)
	}

	// Merging current options with ones has been set in constructor
	options = transport.MergeOptions(n.options, options...)

	// Building request
	req, err := http2.NewRequestWithContext(ctx, method, url, reqData)
	if err != nil {
		return 0, nil, nil, errors.Join(
			errors.New("transport:native: error building request"),
			err,
		)
	}

	// Applying options to request
	for _, o := range options {
		if x, ok := o.(golangNativeRequestOption); ok {
			x.ApplyOnNativeRequest(req)
		}
	}

	// Building client
	cl := &http2.Client{
		Transport: n.transport,
	}

	// Applying options to client
	for _, o := range options {
		if x, ok := o.(golangNativeClientOption); ok {
			x.ApplyOnNativeClient(cl)
		}
	}

	// Performing request
	before := time.Now()
	resp, err := cl.Do(req)
	if err != nil {
		if n.sink != nil {
			n.sink(method, url, req.Header, nil, reqBuf.Bytes(), nil, 0, time.Since(before), err)
		}
		return 0, nil, nil, errors.Join(
			errors.New("transport:native: error performing request"),
			err,
		)
	}

	// Preparing data for log
	var respBuf bytes.Buffer
	var respData io.ReadCloser
	if n.sink != nil {
		respData = &teeOnCompleteReadCloser{
			Reader: io.TeeReader(resp.Body, &respBuf),
			closer: resp.Body,
			onComplete: sync.OnceFunc(func() {
				n.sink(method, url, req.Header, resp.Header, reqBuf.Bytes(), respBuf.Bytes(), resp.StatusCode, time.Since(before), nil)
			}),
		}
	} else {
		respData = resp.Body
	}

	return resp.StatusCode, resp.Header, respData, nil
}

// teeReadCloser wraps rc so that all bytes read from it are also written to w.
func teeReadCloser(rc io.ReadCloser, w io.Writer) io.ReadCloser {
	return struct {
		io.Reader
		io.Closer
	}{
		Reader: io.TeeReader(rc, w),
		Closer: rc,
	}
}

// teeOnCompleteReadCloser calls onComplete exactly once on EOF or Close.
type teeOnCompleteReadCloser struct {
	io.Reader
	closer     io.Closer
	onComplete func()
}

func (t *teeOnCompleteReadCloser) Read(p []byte) (n int, err error) {
	n, err = t.Reader.Read(p)
	if err == io.EOF {
		t.onComplete()
	}
	return
}

func (t *teeOnCompleteReadCloser) Close() error {
	t.onComplete()
	return t.closer.Close()
}
