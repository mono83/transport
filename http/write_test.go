package http

import (
	"io"
	"testing"
)

func TestWriteBytes(t *testing.T) {
	t.Run("returns nil for empty slice", func(t *testing.T) {
		if got := WriteBytes([]byte{}); got != nil {
			t.Fatalf("expected nil, got %v", got)
		}
	})

	t.Run("returns nil for nil slice", func(t *testing.T) {
		if got := WriteBytes(nil); got != nil {
			t.Fatalf("expected nil, got %v", got)
		}
	})

	t.Run("reads back the original bytes", func(t *testing.T) {
		data := []byte("hello")
		rc := WriteBytes(data)
		if rc == nil {
			t.Fatal("unexpected nil")
		}
		defer rc.Close()
		got, err := io.ReadAll(rc)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(got) != "hello" {
			t.Fatalf("got %q, want %q", got, "hello")
		}
	})

	t.Run("close does not error", func(t *testing.T) {
		rc := WriteBytes([]byte("x"))
		if err := rc.Close(); err != nil {
			t.Fatalf("unexpected error on Close: %v", err)
		}
	})

	t.Run("reflects mutations to the original slice", func(t *testing.T) {
		// Documents the no-copy behaviour: WriteBytes retains a reference
		// to the original slice, so mutations are visible to the reader.
		data := []byte("hello")
		rc := WriteBytes(data)
		defer rc.Close()
		data[0] = 'H'
		got, err := io.ReadAll(rc)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(got) != "Hello" {
			t.Fatalf("got %q, want %q — no-copy behaviour changed", got, "Hello")
		}
	})
}

func TestWriteString(t *testing.T) {
	t.Run("returns nil for empty string", func(t *testing.T) {
		if got := WriteString(""); got != nil {
			t.Fatalf("expected nil, got %v", got)
		}
	})

	t.Run("reads back the original string", func(t *testing.T) {
		rc := WriteString("world")
		if rc == nil {
			t.Fatal("unexpected nil")
		}
		defer rc.Close()
		got, err := io.ReadAll(rc)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(got) != "world" {
			t.Fatalf("got %q, want %q", got, "world")
		}
	})

	t.Run("close does not error", func(t *testing.T) {
		rc := WriteString("x")
		if err := rc.Close(); err != nil {
			t.Fatalf("unexpected error on Close: %v", err)
		}
	})
}
