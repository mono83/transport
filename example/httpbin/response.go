package httpbin

type Response struct {
	Headers map[string]string
	Origin  string
	Data    string
	URL     string
}
