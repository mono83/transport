package json

import (
	"errors"
	"io"
	"strings"
	"testing"
)

type testStruct struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func TestReadJSON(t *testing.T) {
	t.Run("decodes valid JSON", func(t *testing.T) {
		rc := io.NopCloser(strings.NewReader(`{"name":"foo","value":42}`))
		got, err := ReadJSON[testStruct](200, nil, rc, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Name != "foo" || got.Value != 42 {
			t.Fatalf("got %+v, want {Name:foo Value:42}", got)
		}
	})

	t.Run("propagates pre-existing error", func(t *testing.T) {
		sentinel := errors.New("transport error")
		rc := io.NopCloser(strings.NewReader(`{"name":"foo"}`))
		got, err := ReadJSON[testStruct](0, nil, rc, sentinel)
		if !errors.Is(err, sentinel) {
			t.Fatalf("got err %v, want %v", err, sentinel)
		}
		if got != nil {
			t.Fatalf("expected nil, got %+v", got)
		}
	})

	t.Run("returns error on invalid JSON", func(t *testing.T) {
		rc := io.NopCloser(strings.NewReader(`not json`))
		got, err := ReadJSON[testStruct](200, nil, rc, nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if got != nil {
			t.Fatalf("expected nil, got %+v", got)
		}
	})

	t.Run("decodes into scalar type", func(t *testing.T) {
		rc := io.NopCloser(strings.NewReader(`"hello"`))
		got, err := ReadJSON[string](200, nil, rc, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if *got != "hello" {
			t.Fatalf("got %q, want %q", *got, "hello")
		}
	})
}
