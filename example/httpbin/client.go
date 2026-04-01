package httpbin

import (
	"context"
	"github.com/mono83/transport/http"
	"io"
)

type Client struct {
	Transport http.Transport
}

func (h Client) ExecuteRequest(
	ctx context.Context,
	method string,
	url string,
	reqData io.ReadCloser,
	options ...any,
) (status int, respHeaders map[string][]string, respData io.ReadCloser, err error) {
	url = http.JoinURL("https://httpbin.org", url)
	return h.Transport.ExecuteRequest(ctx, method, url, reqData, options...)
}
