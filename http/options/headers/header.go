package headers

import "net/http"

// SetHeaderOption is a call option that sets a single HTTP request header,
// replacing any existing values for that header name.
type SetHeaderOption struct {
	Name, Value string
}

// ApplyOnNativeRequest applies the header option to r, replacing any
// previously set values for the same header name.
func (s SetHeaderOption) ApplyOnNativeRequest(r *http.Request) { r.Header.Set(s.Name, s.Value) }

// WithHeader returns a [SetHeaderOption] that sets header n to value v.
func WithHeader(n, v string) SetHeaderOption { return SetHeaderOption{n, v} }

// WithUserAgent returns a [SetHeaderOption] that sets the User-Agent header to v.
func WithUserAgent(v string) SetHeaderOption { return WithHeader("User-Agent", v) }
