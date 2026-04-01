package json

import (
	"io"
	"strings"
	"testing"
)

func TestWrite(t *testing.T) {
	t.Run("encodes struct", func(t *testing.T) {
		rc := Write(testStruct{Name: "foo", Value: 42})
		defer rc.Close()
		got, err := io.ReadAll(rc)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := `{"name":"foo","value":42}`
		if strings.TrimSpace(string(got)) != want {
			t.Fatalf("got %q, want %q", strings.TrimSpace(string(got)), want)
		}
	})

	t.Run("encodes nil", func(t *testing.T) {
		rc := Write(nil)
		defer rc.Close()
		got, err := io.ReadAll(rc)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if strings.TrimSpace(string(got)) != "null" {
			t.Fatalf("got %q, want %q", strings.TrimSpace(string(got)), "null")
		}
	})
}

func TestWriteIndent(t *testing.T) {
	t.Run("encodes struct with indentation", func(t *testing.T) {
		rc := WriteIndent(testStruct{Name: "bar", Value: 7})
		defer rc.Close()
		got, err := io.ReadAll(rc)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := "{\n  \"name\": \"bar\",\n  \"value\": 7\n}"
		if strings.TrimSpace(string(got)) != want {
			t.Fatalf("got %q, want %q", strings.TrimSpace(string(got)), want)
		}
	})
}
