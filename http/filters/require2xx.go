package filters

import (
	"fmt"
	"io"
)

// Require2xx is a filter that passes through the response unchanged when the
// status code is in the 2xx range (200–299), or closes the body and returns
// [ErrExpected2xx] otherwise.
// It is a no-op when err is already set.
//
// It accepts the four return values of [http.Transport.ExecuteRequest] directly,
// so it can be inlined in a call chain:
//
//	json.ReadJSON[T](filters.Require2xx(t.ExecuteRequest(ctx, "POST", url, body)))
func Require2xx(status int, respHeaders map[string][]string, respData io.ReadCloser, err error) (int, map[string][]string, io.ReadCloser, error) {
	if err == nil {
		if status < 200 || status > 299 {
			if respData != nil {
				respData.Close()
			}
			return status, nil, nil, ErrExpected2xx(status)
		}
	}
	return status, respHeaders, respData, err
}

// ErrExpected2xx is returned by [Require2xx] when the response status is outside
// the 2xx range. Its underlying int value is the actual status code received.
type ErrExpected2xx int

func (e ErrExpected2xx) Error() string {
	return fmt.Sprintf("expected HTTP response code 2XX but got %d", e)
}
