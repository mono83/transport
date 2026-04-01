package transport

// MergeOptions appends the variadic options two to the base slice one,
// returning a new combined []any. The original slices are not modified.
func MergeOptions(one []any, two ...any) []any {
	if len(two) == 0 {
		return one
	}
	if len(one) == 0 {
		return two
	}
	out := make([]any, 0, len(one)+len(two))
	out = append(out, one...)
	out = append(out, two...)
	return out
}
