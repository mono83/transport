package xml

import (
	"encoding/xml"
	"github.com/mono83/transport/http"
	"io"
)

// Write encodes obj as XML and returns the result as an [io.ReadCloser].
// Use with headers.WithXMLContentType to set the Content-Type header.
func Write(obj any) io.ReadCloser {
	return http.EncoderPipe(func(w io.WriteCloser) http.Encoder {
		return xml.NewEncoder(w)
	}, obj)
}

// WriteIndent encodes obj as indented XML (two-space indent) and returns
// the result as an [io.ReadCloser].
// Use with headers.WithXMLContentType to set the Content-Type header.
func WriteIndent(obj any) io.ReadCloser {
	return http.EncoderPipe(func(w io.WriteCloser) http.Encoder {
		enc := xml.NewEncoder(w)
		enc.Indent("", "  ")
		return enc
	}, obj)
}
