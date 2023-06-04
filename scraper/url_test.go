package scraper

import (
	"testing"
)

func TestRootDomain(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{"http://example.com/something", "http://example.com"},
		{"http://domain/1/2", "http://domain"},
		{"https://domain/3/4", "https://domain"},
		{"https://example.com/", "https://example.com"},
		{"https://example.com", "https://example.com"},
		{"/route", ""},
		{"/mitsos/mpampis", ""},
	}
	for _, test := range tests {
		if got := rootDomain(test.input); got != test.want {
			t.Errorf(`rootDomain(%q) = %q, want %q`, test.input, got, test.want)
		}
	}
}

func TestTrimAfterHash(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{"/mitsos#fragment", "/mitsos"},
		{"http://example.com/about", "http://example.com/about"},
		{"http://domain.org/something#abc", "http://domain.org/something"},
		{"", ""},
		{"#efg", ""},
	}
	for _, test := range tests {
		if got := trimAfterHash(test.input); got != test.want {
			t.Errorf(`trimAfterHash(%q) = %q, want %q`, test.input, got, test.want)
		}
	}
}

func TestPrefixRoot(t *testing.T) {
	var tests = []struct {
		parent string
		v      string
		want   string
	}{
		{"http://example.com/about/us/", "./abc", "http://example.com/about/abc"},
		{"http://example.com/about/us", "./abc", "http://example.com/about/abc"},
		{"http://example.com/about/", "/efg", "http://example.com/efg"},
		{"http://example.com/about", "/efg", "http://example.com/efg"},
		{"http://example.com/", "http://domain.org", "http://domain.org"},
		{"http://example.com/something/", "abc.txt", "http://example.com/something/abc.txt"},
		{"http://example.com/something", "abc.txt", "http://example.com/abc.txt"},
		{"http://example.com/", "about/", "http://example.com/about"},
		{"http://example.com", "about/", "http://example.com/about"},
	}
	for _, test := range tests {
		if got := prefixRoot(test.parent, test.v); got != test.want {
			t.Errorf(`prefixRoot(%q, %q) = %q, want %q`, test.parent, test.v, got, test.want)
		}
	}
}
