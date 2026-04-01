package http

import "io"

// ReadBytes reads the entire response body and returns it as a byte slice.
// It is designed to be used as a direct pass-through for [Transport.ExecuteRequest]
// return values. The body is always closed before returning when non-nil,
// regardless of whether err is set.
func ReadBytes(
	_ int,
	_ map[string][]string,
	data io.ReadCloser,
	err error,
) ([]byte, error) {
	if data != nil {
		defer data.Close()
	}
	if err != nil {
		return nil, err
	}

	return io.ReadAll(data)
}

// ReadString reads the entire response body and returns it as a string.
// It delegates to [ReadBytes] and converts the result. The body is always
// closed before returning when non-nil, regardless of whether err is set.
func ReadString(
	status int,
	headers map[string][]string,
	data io.ReadCloser,
	err error,
) (string, error) {
	bts, err := ReadBytes(status, headers, data, err)
	if err != nil {
		return "", err
	}
	return string(bts), nil
}
