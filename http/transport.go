package http

import (
	"context"
	"io"
)

// Transport is the low-level HTTP transport contract.
//
// Implementations execute a single HTTP request and return the raw response
// components. The caller is responsible for closing the returned respData body.
// Options are transport-specific values; unknown option types must be silently
// ignored.
type Transport interface {
	ExecuteRequest(
		ctx context.Context,
		method string,
		url string,
		reqData io.ReadCloser,
		options ...any,
	) (
		status int,
		respHeaders map[string][]string,
		respData io.ReadCloser,
		err error,
	)
}
