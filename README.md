# transport

[![CI](https://github.com/mono83/transport/actions/workflows/ci.yml/badge.svg)](https://github.com/mono83/transport/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/mono83/transport.svg)](https://pkg.go.dev/github.com/mono83/transport)
[![Go Report Card](https://goreportcard.com/badge/github.com/mono83/transport)](https://goreportcard.com/report/github.com/mono83/transport)
![Go version](https://img.shields.io/badge/go-1.21+-blue)

> **Experimental** — API may change without notice.

A small, generic HTTP transport abstraction for Go. It separates the concern of *how* a request is executed (the transport layer) from *what* is executed (the call layer), making HTTP clients easy to compose, test, and extend.

## Requirements

Go 1.21 or later.

## Installation

```sh
go get github.com/mono83/transport
```

## Overview

The library is built around two interfaces:

```
Call[R, T]      — a typed operation that uses transport T to produce result R
Transport       — executes a single HTTP request, returns raw response components
```

The root package defines the generic primitives. Sub-packages provide a concrete `net/http` transport, request/response codecs (JSON, XML, form, multipart), helpers for common headers, and optional structured logging.

```
github.com/mono83/transport
├── http/                       Transport interface + read/write helpers + Stub
│   ├── native/                 net/http-backed Transport implementation with transparent gzip/deflate decompression
│   ├── filters/                Status-code guards (Require200, Require2xx)
│   ├── json/                   JSON request/response codec
│   ├── xml/                    XML request/response codec
│   ├── form/                   application/x-www-form-urlencoded codec
│   ├── multipart/              multipart/form-data request builder
│   ├── log/                    Sink type + stdout logger
│   └── options/
│       ├── headers/            Header, auth, and content-type options
│       └── timeout.go          Per-request timeout option
```

## Quick start

A runnable example is in [example/httpbin/](example/httpbin/). It demonstrates a typed client for [httpbin.org](https://httpbin.org):

**[example/httpbin/main/main.go](example/httpbin/main/main.go)** — entry point:
```go
t := native.NewWithLog(log.Stdout)
hb := httpbin.Client{Transport: t}

fmt.Println(httpbin.PostCall{Name: "foo", Value: 99}.Execute(context.Background(), hb))
```

**[example/httpbin/client.go](example/httpbin/client.go)** — wraps the transport with a base URL:
```go
type Client struct {
    Transport http.Transport
}

func (h Client) ExecuteRequest(ctx context.Context, method, url string, reqData io.ReadCloser, options ...any) (int, map[string][]string, io.ReadCloser, error) {
    url = http.JoinURL("https://httpbin.org", url)
    return h.Transport.ExecuteRequest(ctx, method, url, reqData, options...)
}
```

**[example/httpbin/postCall.go](example/httpbin/postCall.go)** — a typed call implementing `Call[*Response, Client]`:
```go
type PostCall struct {
    Name  string
    Value int
}

func (p PostCall) Execute(ctx context.Context, t Client) (*Response, error) {
    return json.ReadJSON[Response](t.ExecuteRequest(ctx, "POST", "/post", json.WriteIndent(p)))
}
```

Run it:
```sh
go run ./example/httpbin/main
```

## Defining typed calls

Implement `Call[R, T]` or use the `CallFunc` adapter to bundle a request into a reusable, testable unit. The pattern shown in [example/httpbin/postCall.go](example/httpbin/postCall.go) applies generally:

```go
// Call is self-contained: method, path, serialisation, and deserialisation
// are all in one place. The client type T is whatever your service wrapper is.
func (c MyCall) Execute(ctx context.Context, t MyClient) (*MyResponse, error) {
    return json.ReadJSON[MyResponse](t.ExecuteRequest(ctx, "POST", "/endpoint", json.Write(c)))
}
```

## Request body helpers

| Helper | Description |
|---|---|
| `http.WriteBytes(b)` | Wraps a `[]byte` as a request body |
| `http.WriteString(s)` | Wraps a `string` as a request body |
| `json.Write(obj)` / `json.WriteIndent(obj)` | JSON-encodes obj into a streaming body |
| `xml.Write(obj)` / `xml.WriteIndent(obj)` | XML-encodes obj into a streaming body |
| `form.Write(values)` | URL-encodes `url.Values` as a request body |
| `multipart.Write(parts, files...)` | Builds a `multipart/form-data` body |

All helpers return `nil` for empty input, which signals "no body" to the transport.

## Response body helpers

| Helper | Description |
|---|---|
| `http.ReadBytes(...)` | Reads the full body as `[]byte` |
| `http.ReadString(...)` | Reads the full body as `string` |
| `json.ReadJSON[T](...)` | Decodes the body as JSON into `*T` |
| `xml.ReadXML[T](...)` | Decodes the body as XML into `*T` |
| `form.Read(...)` | Parses the body as `url.Values` |

All helpers accept the four return values of `ExecuteRequest` directly (spread with `...` or passed inline), close the body, and return a typed result plus an error.

## Filters

The `http/filters` package provides status-code guards that sit between `ExecuteRequest` and a response reader. They accept and return the same four-value signature, so they compose inline:

```go
// Require exactly 200 OK
result, err := json.ReadJSON[T](filters.Require200(t.ExecuteRequest(ctx, "GET", url, nil)))

// Require any 2xx
result, err := json.ReadJSON[T](filters.Require2xx(t.ExecuteRequest(ctx, "POST", url, body)))
```

When the status code does not match, the filter closes the response body and returns a typed error that carries the actual code:

```go
result, err := json.ReadJSON[T](filters.Require200(t.ExecuteRequest(ctx, "GET", url, nil)))
if e, ok := err.(filters.ErrExpected200); ok {
    fmt.Println("unexpected status:", int(e))
}
```

Filters are a no-op when a transport-level error is already present.

| Filter | Passes when |
|---|---|
| `filters.Require200` | status == 200 |
| `filters.Require2xx` | 200 ≤ status ≤ 299 |

## Options

Options are plain values passed as variadic `...any` to `ExecuteRequest`. Unknown option types are silently ignored.

```go
// Set a header
headers.WithHeader("X-Request-ID", "abc123")
headers.WithBearerToken("my-token")
headers.WithBasicAuth("user", "pass")
headers.WithAPIKey("X-Api-Key", "key")

// Content negotiation
headers.WithJSONContentType()
headers.WithAcceptJSON()
headers.WithXMLContentType()
headers.WithAcceptXML()

// Compression — pairs with the native transport's transparent decompression
headers.WithAcceptEncodingGzipDeflate()

// Per-request timeout (native transport)
options.Timeout(5 * time.Second)
```

Constructor-level options (passed to `native.New`) are merged with per-call options at request time.

## Transparent decompression

The native transport automatically decompresses response bodies based on the `Content-Encoding` response header. Decompression is invisible to callers — response body helpers (`ReadJSON`, `ReadBytes`, etc.) always receive plain text.

| `Content-Encoding` | Handling |
|---|---|
| `gzip` | Decompressed via `compress/gzip` |
| `deflate` | Auto-detects zlib framing (RFC 1950) or raw DEFLATE (RFC 1951) |
| anything else | Passed through unchanged |

To advertise compression support to the server, pass `WithAcceptEncodingGzipDeflate` as an option:

```go
result, err := json.ReadJSON[T](t.ExecuteRequest(ctx, "GET", url, nil,
    headers.WithAcceptEncodingGzipDeflate()))
```

Or set it once at construction time so every request from that transport includes it:

```go
t := native.New(headers.WithAcceptEncodingGzipDeflate())
```

Decompression happens before logging, so the log sink always receives decoded bytes.

## Logging

`log.Sink` is a function type that receives the full record of an exchange — method, URL, request/response headers, bodies, status, elapsed time, and any error — after the response body is consumed or closed.

```go
// Built-in verbose stdout sink
t := native.NewWithLog(log.Stdout)

// Custom sink
var sink log.Sink = func(method, url string, reqH, respH map[string][]string,
    reqBody, respBody []byte, status int, elapsed time.Duration, err error) {
    slog.Info("http", "method", method, "url", url, "status", status, "elapsed", elapsed)
}
t := native.NewWithLog(sink)
```

## Multipart uploads

```go
form := multipart.Write(
    []multipart.Part{
        {Name: "author", Value: "alice"},
    },
    multipart.FilePart{
        Name:     "file",
        FileName: "report.csv",
        Body:     f, // any io.ReadCloser; closed automatically after copy
    },
)

t.ExecuteRequest(ctx, "POST", url, form,
    headers.WithHeader("Content-Type", form.ContentType()))
```

## Testing

Three ready-made `Transport` implementations are available for unit tests. All are safe to reuse across multiple calls.

`http.Stub` — full control over status, headers, and body:

```go
stub := http.Stub{
    Status:       200,
    ResponseData: []byte(`{"id":1}`),
}
result, err := json.ReadJSON[MyType](stub.ExecuteRequest(ctx, "GET", "/anything", nil))
```

Simulate a transport-level error:

```go
stub := http.Stub{Error: errors.New("connection refused")}
_, err := json.ReadJSON[MyType](stub.ExecuteRequest(ctx, "GET", "/anything", nil))
// err != nil
```

`http.StubBytes` and `http.StubString` — lightweight stubs for when only the body matters:

```go

result, err := json.ReadJSON[MyType](http.StubBytes(`{"id":1}`).ExecuteRequest(ctx, "GET", "/", nil))
result, err := json.ReadJSON[MyType](http.StubString(`{"id":1}`).ExecuteRequest(ctx, "GET", "/", nil))
```
