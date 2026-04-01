package form

import (
	"io"
	"net/url"

	"github.com/mono83/transport/http"
)

// Write encodes values as an application/x-www-form-urlencoded body
// (e.g. "foo=bar&baz=3") and returns it as an [io.ReadCloser].
// Returns nil for nil or empty values, signalling "no body" to the transport.
// Use with headers.WithFormURLEncodedContentType to set the Content-Type header.
func Write(values url.Values) io.ReadCloser {
	if len(values) == 0 {
		return nil
	}
	return http.WriteString(values.Encode())
}
