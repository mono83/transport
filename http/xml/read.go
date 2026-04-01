package xml

import (
	"encoding/xml"
	"io"
)

// ReadXML decodes the response body as XML into a value of type T and
// returns a pointer to it. The body is always closed before returning when
// non-nil, regardless of whether err is set. Returns an error if the body
// cannot be decoded as valid XML.
func ReadXML[T any](
	_ int,
	_ map[string][]string,
	data io.ReadCloser,
	err error,
) (*T, error) {
	if data != nil {
		defer data.Close()
	}
	if err != nil {
		return nil, err
	}

	var t T
	if err := xml.NewDecoder(data).Decode(&t); err != nil {
		return nil, err
	}

	return &t, nil
}
