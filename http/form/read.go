package form

import (
	"io"
	"net/url"

	"github.com/mono83/transport/http"
)

// Read reads the response body as an application/x-www-form-urlencoded payload
// and returns the parsed [url.Values].
// The body is always closed before returning when non-nil, regardless of
// whether err is set.
func Read(
	_ int,
	_ map[string][]string,
	data io.ReadCloser,
	err error,
) (url.Values, error) {
	b, err := http.ReadBytes(0, nil, data, err)
	if err != nil {
		return nil, err
	}
	return url.ParseQuery(string(b))
}
