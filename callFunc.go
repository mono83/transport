package transport

import "context"

// CallFunc is a function type that implements [Call] for cases where defining a
// named struct is unnecessary.
//
// Any function with the matching signature can be converted to CallFunc and
// passed wherever a Call[R, T] is expected.
type CallFunc[R any, T any] func(ctx context.Context, transport T) (R, error)

// Execute calls the underlying function and returns its result.
func (c CallFunc[R, T]) Execute(ctx context.Context, transport T) (R, error) {
	return c(ctx, transport)
}
