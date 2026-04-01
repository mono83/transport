package http

import (
	"io"
)

// Encoder is implemented by types that can encode a value into a stream
// (e.g. [encoding/json.Encoder], [encoding/xml.Encoder]).
type Encoder interface {
	Encode(any) error
}

// EncoderPipe encodes obj using the Encoder produced by buildEncoder and
// returns the encoded bytes as an [io.ReadCloser]. Encoding runs in a
// goroutine; if it fails, the error is propagated to the reader via
// [io.PipeWriter.CloseWithError].
func EncoderPipe(buildEncoder func(io.WriteCloser) Encoder, obj any) io.ReadCloser {
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()

		if err := buildEncoder(pw).Encode(obj); err != nil {
			pw.CloseWithError(err)
		}
	}()
	return pr
}
