package log

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"
)

// Stdout is a [Sink] that prints a verbose dump of each request and response
// to stdout, including headers and bodies rendered as strings.
var Stdout Sink = stdoutSink

var outputCounter uint64

func stdoutSink(
	method string,
	url string,
	reqHeaders map[string][]string,
	respHeaders map[string][]string,
	reqData []byte,
	respData []byte,
	status int,
	elapsed time.Duration,
	err error,
) {
	var sb strings.Builder

	fmt.Fprintf(&sb, "--- [%d:%s:%d] %s elapsed=%s\n", atomic.AddUint64(&outputCounter, 1), method, status, url, elapsed.Round(time.Millisecond))

	writeHeaders(&sb, "req headers", reqHeaders)
	if len(reqData) > 0 {
		fmt.Fprintf(&sb, "req body:\n%s\n", reqData)
	}

	if err != nil {
		fmt.Fprintf(&sb, "error: %v\n", err)
	} else {
		writeHeaders(&sb, "resp headers", respHeaders)
		if len(respData) > 0 {
			fmt.Fprintf(&sb, "resp body:\n%s\n", respData)
		}
	}

	fmt.Print(sb.String())
}

func writeHeaders(sb *strings.Builder, label string, headers map[string][]string) {
	if len(headers) == 0 {
		return
	}
	fmt.Fprintf(sb, "%s:\n", label)
	for name, values := range headers {
		for _, v := range values {
			fmt.Fprintf(sb, "  %s: %s\n", name, v)
		}
	}
}
