package util

import (
	"net/http"
	"testing"
	"time"
)

func TestNormalizeURLPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{"empty string", "", "/"},
		{"no leading slash", "foo", "/foo"},
		{"has leading slash", "/foo", "/foo"},
		{"multiple segments no leading slash", "foo/bar", "/foo/bar"},
		{"multiple segments with leading slash", "/foo/bar", "/foo/bar"},
		{"just a slash", "/", "/"},
		{"double leading slash", "//foo", "//foo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeURLPath(tt.path)
			if got != tt.want {
				t.Errorf("NormalizeURLPath(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestNormalizeBaseURL(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		want    string
	}{
		{"empty string", "", ""},
		{"no trailing slash", "http://example.com", "http://example.com"},
		{"single trailing slash", "http://example.com/", "http://example.com"},
		{"multiple trailing slashes", "http://example.com///", "http://example.com"},
		{"path with trailing slash", "http://example.com/v1/", "http://example.com/v1"},
		{"just a scheme and host", "https://host", "https://host"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeBaseURL(tt.baseURL)
			if got != tt.want {
				t.Errorf("NormalizeBaseURL(%q) = %q, want %q", tt.baseURL, got, tt.want)
			}
		})
	}
}

func TestParseRetryAfter(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		wantZero bool
	}{
		{name: "integer seconds", value: "120"},
		{name: "zero seconds", value: "0", wantZero: true},
		{name: "HTTP-date future", value: time.Now().Add(10 * time.Second).UTC().Format(http.TimeFormat)},
		{name: "HTTP-date past", value: time.Now().Add(-10 * time.Second).UTC().Format(http.TimeFormat), wantZero: true},
		{name: "negative integer", value: "-5", wantZero: true},
		{name: "invalid value", value: "abc", wantZero: true},
		{name: "empty string", value: "", wantZero: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := ParseRetryAfter(tt.value)
			if tt.wantZero && d != 0 {
				t.Errorf("expected zero duration, got %v", d)
			}
			if !tt.wantZero && d <= 0 {
				t.Errorf("expected positive duration, got %v", d)
			}
		})
	}
}
