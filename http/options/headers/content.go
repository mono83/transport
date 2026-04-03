package headers

// Content-Type values for common media types.
const (
	contentTypeJSON           = "application/json"
	contentTypeXML            = "application/xml"
	contentTypeFormURLEncoded = "application/x-www-form-urlencoded"
	contentTypeMultipartForm  = "multipart/form-data"
	contentTypePlainText      = "text/plain; charset=utf-8"
	contentTypeHTML           = "text/html; charset=utf-8"
	contentTypeOctetStream    = "application/octet-stream"
)

// WithContentType returns a [SetHeaderOption] that sets the Content-Type header.
func WithContentType(v string) SetHeaderOption { return WithHeader("Content-Type", v) }

// WithAccept returns a [SetHeaderOption] that sets the Accept header.
func WithAccept(v string) SetHeaderOption { return WithHeader("Accept", v) }

// WithJSONContentType returns a [SetHeaderOption] that sets Content-Type to application/json.
func WithJSONContentType() SetHeaderOption { return WithContentType(contentTypeJSON) }

// WithXMLContentType returns a [SetHeaderOption] that sets Content-Type to application/xml.
func WithXMLContentType() SetHeaderOption { return WithContentType(contentTypeXML) }

// WithAcceptJSON returns a [SetHeaderOption] that sets Accept to application/json.
func WithAcceptJSON() SetHeaderOption { return WithAccept(contentTypeJSON) }

// WithAcceptXML returns a [SetHeaderOption] that sets Accept to application/xml.
func WithAcceptXML() SetHeaderOption { return WithAccept(contentTypeXML) }

// WithFormURLEncodedContentType returns a [SetHeaderOption] that sets Content-Type to application/x-www-form-urlencoded.
func WithFormURLEncodedContentType() SetHeaderOption {
	return WithContentType(contentTypeFormURLEncoded)
}

// WithMultipartFormContentType returns a [SetHeaderOption] that sets Content-Type to multipart/form-data.
func WithMultipartFormContentType() SetHeaderOption { return WithContentType(contentTypeMultipartForm) }

// WithPlainTextContentType returns a [SetHeaderOption] that sets Content-Type to text/plain; charset=utf-8.
func WithPlainTextContentType() SetHeaderOption { return WithContentType(contentTypePlainText) }

// WithHTMLContentType returns a [SetHeaderOption] that sets Content-Type to text/html; charset=utf-8.
func WithHTMLContentType() SetHeaderOption { return WithContentType(contentTypeHTML) }

// WithOctetStreamContentType returns a [SetHeaderOption] that sets Content-Type to application/octet-stream.
func WithOctetStreamContentType() SetHeaderOption { return WithContentType(contentTypeOctetStream) }

// WithAcceptEncodingGzipDeflate returns a [SetHeaderOption] that sets Accept-Encoding to "gzip, deflate",
// signalling that the caller can handle compressed responses. The native transport decompresses
// gzip and deflate responses transparently, so this option pairs naturally with it.
func WithAcceptEncodingGzipDeflate() SetHeaderOption {
	return WithHeader("Accept-Encoding", "gzip, deflate")
}
