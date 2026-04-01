package xml

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
		want := `<testStruct><name>foo</name><value>42</value></testStruct>`
		if strings.TrimSpace(string(got)) != want {
			t.Fatalf("got %q, want %q", strings.TrimSpace(string(got)), want)
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
		want := "<testStruct>\n  <name>bar</name>\n  <value>7</value>\n</testStruct>"
		if strings.TrimSpace(string(got)) != want {
			t.Fatalf("got %q, want %q", strings.TrimSpace(string(got)), want)
		}
	})
}
