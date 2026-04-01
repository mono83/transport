package headers

import "encoding/base64"

// WithBearerToken returns a [SetHeaderOption] that sets the Authorization header
// to "Bearer <token>".
func WithBearerToken(token string) SetHeaderOption {
	return WithHeader("Authorization", "Bearer "+token)
}

// WithBasicAuth returns a [SetHeaderOption] that sets the Authorization header
// using HTTP Basic authentication with the given username and password.
func WithBasicAuth(username, password string) SetHeaderOption {
	credentials := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	return WithHeader("Authorization", "Basic "+credentials)
}

// WithAPIKey returns a [SetHeaderOption] that sets a custom header to the given
// API key value. Common header names are "X-Api-Key" and "X-Auth-Token".
func WithAPIKey(header, key string) SetHeaderOption { return WithHeader(header, key) }
