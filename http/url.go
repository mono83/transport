package http

import "strings"

// JoinURL joins a base URL prefix (e.g. "https://example.com") with a request
// URI (e.g. "/path"), returning a single URL with exactly one slash at the
// join point. Trailing slashes on prefix and leading slashes on uri are
// normalised so the result never contains a double slash at the boundary.
func JoinURL(prefix, uri string) string {
	return strings.TrimRight(prefix, "/") + "/" + strings.TrimLeft(uri, "/")
}
