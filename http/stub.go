package http

import (
	"bytes"
	"context"
	"io"
)

// Stub is a [Transport] implementation for use in unit tests. It returns a
// fixed response for every call to [Stub.ExecuteRequest], regardless of method,
// URL, or request body. Each call creates a fresh reader over ResponseData, so
// the same Stub can be reused across multiple calls.
//
// Set Error to simulate a transport-level failure; Status and ResponseData are
// ignored in that case.
type Stub struct {
	Status          int
	ResponseHeaders map[string][]string
	ResponseData    []byte
	Error           error
}

// ExecuteRequest returns the configured fixed response.
// It ignores ctx, method, url, reqData, and options.
func (s Stub) ExecuteRequest(
	ctx context.Context,
	method string,
	url string,
	reqData io.ReadCloser,
	options ...any,
) (int, map[string][]string, io.ReadCloser, error) {
	if s.Error != nil {
		return 0, nil, nil, s.Error
	}

	var body io.ReadCloser
	if len(s.ResponseData) > 0 {
		body = io.NopCloser(bytes.NewReader(s.ResponseData))
	}

	return s.Status, s.ResponseHeaders, body, nil
}
