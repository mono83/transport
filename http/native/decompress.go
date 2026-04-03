package native

import (
	"bufio"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"errors"
	"io"
	"strings"
)

// decompressBody transparently decompresses respData based on the
// Content-Encoding response header. It handles "gzip" and "deflate";
// all other encodings are passed through unchanged.
func decompressBody(respHeaders map[string][]string, respData io.ReadCloser) (io.ReadCloser, error) {
	encoding := contentEncoding(respHeaders)

	switch encoding {
	case "gzip":
		gr, err := gzip.NewReader(respData)
		if err != nil {
			respData.Close()
			return nil, errors.Join(
				errors.New("transport:native: failed to initialise gzip reader"),
				err,
			)
		}
		return &gzipReadCloser{Reader: gr, underlying: respData}, nil

	case "deflate":
		rc, err := newDeflateReader(respData)
		if err != nil {
			respData.Close()
			return nil, errors.Join(
				errors.New("transport:native: failed to initialise deflate reader"),
				err,
			)
		}
		return rc, nil

	default:
		return respData, nil
	}
}

// contentEncoding returns the lowercased Content-Encoding value from headers,
// or an empty string if absent.
func contentEncoding(headers map[string][]string) string {
	for k, v := range headers {
		if strings.EqualFold(k, "Content-Encoding") && len(v) > 0 {
			return strings.ToLower(strings.TrimSpace(v[0]))
		}
	}
	return ""
}

// newDeflateReader auto-detects whether the stream uses zlib framing (RFC 1950)
// or raw DEFLATE (RFC 1951, used by many real-world servers).
func newDeflateReader(rc io.ReadCloser) (io.ReadCloser, error) {
	br := bufio.NewReader(rc)

	b, err := br.Peek(2)
	if err != nil && len(b) < 2 {
		return nil, err
	}

	// zlib streams start with 0x78 and a second byte that makes the
	// two-byte big-endian value divisible by 31 (CMF/FLG check).
	if len(b) == 2 && b[0] == 0x78 && isZlibFLG(b[1]) {
		zr, zerr := zlib.NewReader(br)
		if zerr != nil {
			return nil, zerr
		}
		return &zlibReadCloser{ReadCloser: zr, underlying: rc}, nil
	}

	// Fall back to raw DEFLATE.
	return &flateReadCloser{ReadCloser: flate.NewReader(br), underlying: rc}, nil
}

// isZlibFLG reports whether cmf=0x78 combined with flg produces a header
// value (0x78<<8 | flg) that is divisible by 31, as required by RFC 1950.
func isZlibFLG(flg byte) bool {
	return (0x7800|uint16(flg))%31 == 0
}

type gzipReadCloser struct {
	*gzip.Reader
	underlying io.Closer
}

func (g *gzipReadCloser) Close() error {
	return errors.Join(g.Reader.Close(), g.underlying.Close())
}

type zlibReadCloser struct {
	io.ReadCloser
	underlying io.Closer
}

func (z *zlibReadCloser) Close() error {
	return errors.Join(z.ReadCloser.Close(), z.underlying.Close())
}

type flateReadCloser struct {
	io.ReadCloser
	underlying io.Closer
}

func (f *flateReadCloser) Close() error {
	return errors.Join(f.ReadCloser.Close(), f.underlying.Close())
}
