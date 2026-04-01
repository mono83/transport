package xml

import (
	"errors"
	"io"
	"strings"
	"testing"
)

type testStruct struct {
	Name  string `xml:"name"`
	Value int    `xml:"value"`
}

func TestReadXML(t *testing.T) {
	t.Run("decodes valid XML", func(t *testing.T) {
		body := `<testStruct><name>foo</name><value>42</value></testStruct>`
		rc := io.NopCloser(strings.NewReader(body))
		got, err := ReadXML[testStruct](200, nil, rc, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Name != "foo" || got.Value != 42 {
			t.Fatalf("got %+v, want {Name:foo Value:42}", got)
		}
	})

	t.Run("propagates pre-existing error", func(t *testing.T) {
		sentinel := errors.New("transport error")
		rc := io.NopCloser(strings.NewReader(`<testStruct><name>foo</name></testStruct>`))
		got, err := ReadXML[testStruct](0, nil, rc, sentinel)
		if !errors.Is(err, sentinel) {
			t.Fatalf("got err %v, want %v", err, sentinel)
		}
		if got != nil {
			t.Fatalf("expected nil, got %+v", got)
		}
	})

	t.Run("returns error on invalid XML", func(t *testing.T) {
		rc := io.NopCloser(strings.NewReader(`not xml`))
		got, err := ReadXML[testStruct](200, nil, rc, nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if got != nil {
			t.Fatalf("expected nil, got %+v", got)
		}
	})
}
