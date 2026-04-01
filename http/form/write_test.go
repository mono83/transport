package form

import (
	"io"
	"net/url"
	"testing"
)

func TestWrite(t *testing.T) {
	t.Run("returns nil for nil values", func(t *testing.T) {
		if got := Write(nil); got != nil {
			t.Fatalf("expected nil, got %v", got)
		}
	})

	t.Run("returns nil for empty values", func(t *testing.T) {
		if got := Write(url.Values{}); got != nil {
			t.Fatalf("expected nil, got %v", got)
		}
	})

	t.Run("encodes single field", func(t *testing.T) {
		rc := Write(url.Values{"foo": {"bar"}})
		if rc == nil {
			t.Fatal("unexpected nil")
		}
		defer rc.Close()
		got, err := io.ReadAll(rc)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(got) != "foo=bar" {
			t.Fatalf("got %q, want %q", got, "foo=bar")
		}
	})

	t.Run("encodes multiple fields sorted by key", func(t *testing.T) {
		rc := Write(url.Values{"b": {"2"}, "a": {"1"}})
		defer rc.Close()
		got, err := io.ReadAll(rc)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// url.Values.Encode sorts keys alphabetically
		if string(got) != "a=1&b=2" {
			t.Fatalf("got %q, want %q", got, "a=1&b=2")
		}
	})

	t.Run("encodes multi-value field", func(t *testing.T) {
		rc := Write(url.Values{"x": {"1", "2"}})
		defer rc.Close()
		got, err := io.ReadAll(rc)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(got) != "x=1&x=2" {
			t.Fatalf("got %q, want %q", got, "x=1&x=2")
		}
	})

	t.Run("percent-encodes special characters", func(t *testing.T) {
		rc := Write(url.Values{"q": {"hello world"}})
		defer rc.Close()
		got, err := io.ReadAll(rc)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(got) != "q=hello+world" {
			t.Fatalf("got %q, want %q", got, "q=hello+world")
		}
	})
}
