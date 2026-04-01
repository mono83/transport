package http

import "testing"

func TestJoinURL(t *testing.T) {
	tests := []struct {
		prefix, uri, want string
	}{
		{"https://example.com", "/path", "https://example.com/path"},
		{"https://example.com/", "/path", "https://example.com/path"},
		{"https://example.com", "path", "https://example.com/path"},
		{"https://example.com/", "path", "https://example.com/path"},
		{"https://example.com/base", "/path", "https://example.com/base/path"},
		{"https://example.com/base/", "/path", "https://example.com/base/path"},
		{"https://example.com", "", "https://example.com/"},
		{"https://example.com", "/", "https://example.com/"},
	}

	for _, tt := range tests {
		got := JoinURL(tt.prefix, tt.uri)
		if got != tt.want {
			t.Errorf("JoinURL(%q, %q) = %q, want %q", tt.prefix, tt.uri, got, tt.want)
		}
	}
}
