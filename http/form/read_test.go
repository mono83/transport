package form

import (
	"errors"
	"io"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	t.Run("parses single field", func(t *testing.T) {
		rc := io.NopCloser(strings.NewReader("foo=bar"))
		got, err := Read(200, nil, rc, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Get("foo") != "bar" {
			t.Fatalf("foo: got %q, want %q", got.Get("foo"), "bar")
		}
	})

	t.Run("parses multiple fields", func(t *testing.T) {
		rc := io.NopCloser(strings.NewReader("a=1&b=2"))
		got, err := Read(200, nil, rc, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Get("a") != "1" {
			t.Fatalf("a: got %q, want %q", got.Get("a"), "1")
		}
		if got.Get("b") != "2" {
			t.Fatalf("b: got %q, want %q", got.Get("b"), "2")
		}
	})

	t.Run("parses multi-value field", func(t *testing.T) {
		rc := io.NopCloser(strings.NewReader("x=1&x=2"))
		got, err := Read(200, nil, rc, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got["x"]) != 2 || got["x"][0] != "1" || got["x"][1] != "2" {
			t.Fatalf("x: got %v, want [1 2]", got["x"])
		}
	})

	t.Run("decodes percent-encoded values", func(t *testing.T) {
		rc := io.NopCloser(strings.NewReader("q=hello+world"))
		got, err := Read(200, nil, rc, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Get("q") != "hello world" {
			t.Fatalf("q: got %q, want %q", got.Get("q"), "hello world")
		}
	})

	t.Run("propagates pre-existing error", func(t *testing.T) {
		sentinel := errors.New("transport error")
		rc := io.NopCloser(strings.NewReader("foo=bar"))
		got, err := Read(0, nil, rc, sentinel)
		if !errors.Is(err, sentinel) {
			t.Fatalf("got err %v, want %v", err, sentinel)
		}
		if got != nil {
			t.Fatalf("expected nil values, got %v", got)
		}
	})
}
