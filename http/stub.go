package http

import (
	"bytes"
	"context"
	"io"
	"strings"
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
	_ context.Context,
	_ string,
	_ string,
	_ io.ReadCloser,
	_ ...any,
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

// StubBytes is a [Transport] implementation for use in unit tests. It returns
// a fixed response body from a raw byte slice for every call to
// [StubBytes.ExecuteRequest], with status 0 and no headers.
// Each call creates a fresh reader, so the same StubBytes can be reused across
// multiple calls.
type StubBytes []byte

// ExecuteRequest returns status 0, nil headers, and a fresh reader over the
// byte slice. It ignores ctx, method, url, reqData, and options.
func (s StubBytes) ExecuteRequest(_ context.Context, _ string, _ string, _ io.ReadCloser, _ ...any) (int, map[string][]string, io.ReadCloser, error) {
	return 0, nil, io.NopCloser(bytes.NewReader(s)), nil
}

// StubString is a [Transport] implementation for use in unit tests. It returns
// a fixed response body from a string for every call to
// [StubString.ExecuteRequest], with status 0 and no headers.
// Each call creates a fresh reader, so the same StubString can be reused across
// multiple calls.
type StubString string

// ExecuteRequest returns status 0, nil headers, and a fresh reader over the
// string. It ignores ctx, method, url, reqData, and options.
func (s StubString) ExecuteRequest(_ context.Context, _ string, _ string, _ io.ReadCloser, _ ...any) (int, map[string][]string, io.ReadCloser, error) {
	return 0, nil, io.NopCloser(strings.NewReader(string(s))), nil
}
