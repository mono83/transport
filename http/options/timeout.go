// Package options provides call options for the native HTTP transport.
package options

import (
	"net/http"
	"time"
)

// Timeout is a per-request time limit passed as a call option.
// It sets [net/http.Client.Timeout], which covers the entire request lifecycle
// including redirect follows and reading the response body.
// A zero value means no timeout.
type Timeout time.Duration

func (t Timeout) ApplyOnNativeClient(c *http.Client) {
	c.Timeout = time.Duration(t)
}
