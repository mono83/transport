package http

import (
	"errors"
	"io"
	"strings"
	"testing"
)

func TestReadBytes(t *testing.T) {
	t.Run("reads body and closes it", func(t *testing.T) {
		rc := io.NopCloser(strings.NewReader("hello"))
		got, err := ReadBytes(200, nil, rc, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(got) != "hello" {
			t.Fatalf("got %q, want %q", got, "hello")
		}
	})

	t.Run("returns error without reading body", func(t *testing.T) {
		sentinel := errors.New("transport error")
		rc := io.NopCloser(strings.NewReader("should not be read"))
		got, err := ReadBytes(0, nil, rc, sentinel)
		if !errors.Is(err, sentinel) {
			t.Fatalf("got err %v, want %v", err, sentinel)
		}
		if got != nil {
			t.Fatalf("expected nil bytes, got %v", got)
		}
	})

	t.Run("empty body", func(t *testing.T) {
		rc := io.NopCloser(strings.NewReader(""))
		got, err := ReadBytes(204, nil, rc, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 0 {
			t.Fatalf("expected empty slice, got %v", got)
		}
	})
}

func TestReadString(t *testing.T) {
	t.Run("reads body as string", func(t *testing.T) {
		rc := io.NopCloser(strings.NewReader("world"))
		got, err := ReadString(200, nil, rc, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "world" {
			t.Fatalf("got %q, want %q", got, "world")
		}
	})

	t.Run("returns error without reading body", func(t *testing.T) {
		sentinel := errors.New("transport error")
		rc := io.NopCloser(strings.NewReader("should not be read"))
		got, err := ReadString(0, nil, rc, sentinel)
		if !errors.Is(err, sentinel) {
			t.Fatalf("got err %v, want %v", err, sentinel)
		}
		if got != "" {
			t.Fatalf("expected empty string, got %q", got)
		}
	})
}
