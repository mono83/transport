package log

import "time"

// Sink receives the full record of an HTTP exchange after it completes.
// elapsed is the time from just before ExecuteRequest was called until the
// response body was fully consumed or closed.
// On error, respHeaders and respData are nil.
type Sink func(
	method string,
	url string,
	reqHeaders, respHeaders map[string][]string,
	reqData, respData []byte,
	status int,
	elapsed time.Duration,
	err error,
)
