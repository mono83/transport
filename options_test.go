package transport

import (
	"slices"
	"testing"
)

func TestMergeOptions(t *testing.T) {
	t.Run("merges two non-empty slices", func(t *testing.T) {
		one := []any{"a", "b"}
		got := MergeOptions(one, "c", "d")
		want := []any{"a", "b", "c", "d"}
		if !slices.Equal(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})

	t.Run("returns one when two is empty", func(t *testing.T) {
		one := []any{"a", "b"}
		got := MergeOptions(one)
		if &got[0] != &one[0] {
			t.Fatal("expected same underlying slice as one")
		}
	})

	t.Run("returns two when one is empty", func(t *testing.T) {
		got := MergeOptions(nil, "a", "b")
		want := []any{"a", "b"}
		if !slices.Equal(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})

	t.Run("both empty returns empty", func(t *testing.T) {
		got := MergeOptions(nil)
		if len(got) != 0 {
			t.Fatalf("expected empty, got %v", got)
		}
	})

	t.Run("does not modify original slices", func(t *testing.T) {
		one := []any{"a", "b"}
		oneCopy := []any{"a", "b"}
		MergeOptions(one, "c")
		if !slices.Equal(one, oneCopy) {
			t.Fatal("one was modified")
		}
	})
}
