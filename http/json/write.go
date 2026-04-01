package json

import (
	"encoding/json"
	"github.com/mono83/transport/http"
	"io"
)

// Write encodes obj as JSON and returns the result as an [io.ReadCloser].
// Use with headers.WithJSONContentType to set the Content-Type header.
func Write(obj any) io.ReadCloser {
	return http.EncoderPipe(func(w io.WriteCloser) http.Encoder {
		return json.NewEncoder(w)
	}, obj)
}

// WriteIndent encodes obj as indented JSON (two-space indent) and returns
// the result as an [io.ReadCloser].
// Use with headers.WithJSONContentType to set the Content-Type header.
func WriteIndent(obj any) io.ReadCloser {
	return http.EncoderPipe(func(w io.WriteCloser) http.Encoder {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc
	}, obj)
}
