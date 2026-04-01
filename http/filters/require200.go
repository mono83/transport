package filters

import (
	"fmt"
	"io"
)

// Require200 is a filter that passes through the response unchanged when the
// status code is 200, or closes the body and returns [ErrExpected200] otherwise.
// It is a no-op when err is already set.
//
// It accepts the four return values of [http.Transport.ExecuteRequest] directly,
// so it can be inlined in a call chain:
//
//	json.ReadJSON[T](filters.Require200(t.ExecuteRequest(ctx, "GET", url, nil)))
func Require200(status int, respHeaders map[string][]string, respData io.ReadCloser, err error) (int, map[string][]string, io.ReadCloser, error) {
	if err == nil {
		if status != 200 {
			if respData != nil {
				respData.Close()
			}
			return status, nil, nil, ErrExpected200(status)
		}
	}
	return status, respHeaders, respData, err
}

// ErrExpected200 is returned by [Require200] when the response status is not
// 200. Its underlying int value is the actual status code received.
type ErrExpected200 int

func (e ErrExpected200) Error() string {
	return fmt.Sprintf("expected HTTP response code 200 but got %d", e)
}
