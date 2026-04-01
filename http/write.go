package http

import (
	"bytes"
	"io"
	"strings"
)

// WriteBytes wraps b in an [io.ReadCloser] suitable for use as a request body.
// Returns nil for a nil or empty slice, signalling "no body" to the transport.
// The returned reader retains a reference to b — mutations to b after this
// call are visible to the reader.
func WriteBytes(b []byte) io.ReadCloser {
	if len(b) == 0 {
		return nil
	}

	return io.NopCloser(bytes.NewReader(b))
}

// WriteString wraps s in an [io.ReadCloser] suitable for use as a request body.
// Returns nil for an empty string, signalling "no body" to the transport.
func WriteString(s string) io.ReadCloser {
	if len(s) == 0 {
		return nil
	}

	return io.NopCloser(strings.NewReader(s))
}
