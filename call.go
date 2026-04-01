package transport

import "context"

// Call represents a single executable operation against transport T that
// produces a result of type R.
//
// R is the result type; T is the transport type (e.g. an HTTP client interface).
// Implementations must honour context cancellation and propagate transport errors.
type Call[R any, T any] interface {
	// Execute executes call using given Transport and returns result
	Execute(ctx context.Context, transport T) (R, error)
}
