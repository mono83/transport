// Package native provides a [http.Transport] implementation backed by [net/http].
package native

import (
	"github.com/mono83/transport/http"
	"github.com/mono83/transport/http/log"
	http2 "net/http"
	"time"
)

// New returns an [http.Transport] backed by [net/http] with sensible connection
// pool defaults (100 max idle connections, 90 s idle timeout, 20 per host).
// Options are stored and merged with per-call options at request time.
func New(options ...any) http.Transport {
	return NewWithLog(nil, options...)
}

// NewWithLog is like [New] but records every exchange via s.
// Pass nil to disable logging.
func NewWithLog(s log.Sink, options ...any) http.Transport {
	t := &http2.Transport{
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		MaxIdleConnsPerHost: 20,
	}

	// Applying transport-level options
	for _, o := range options {
		if x, ok := o.(golangNativeTransportOption); ok {
			x.ApplyOnNativeTransport(t)
		}
	}

	return &nativeTransport{
		transport: t,
		sink:      s,
		options:   options,
	}
}
