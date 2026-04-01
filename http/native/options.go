package native

import "net/http"

// golangNativeTransportOption is implemented by options that configure the
// underlying [net/http.Transport] (e.g. TLS settings, proxy, connection limits).
type golangNativeTransportOption interface {
	ApplyOnNativeTransport(*http.Transport)
}

// golangNativeClientOption is implemented by options that configure the
// [net/http.Client] used for a single request (e.g. timeout, redirect policy).
type golangNativeClientOption interface {
	ApplyOnNativeClient(*http.Client)
}

// golangNativeRequestOption is implemented by options that mutate the
// [net/http.Request] before it is sent (e.g. headers, authentication).
type golangNativeRequestOption interface {
	ApplyOnNativeRequest(*http.Request)
}
