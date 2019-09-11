package sitemap

import "testing"

func TestHref(t *testing.T) {
	tests := []struct {
		host, href, expected string
	}{
		{"http://wow.com", "/hey", "http://wow.com/hey"},
		{"http://wow.com/", "https://hey", "https://hey"},
	}
	for i, test := range tests {
		if url := href(test.href, test.host); url != test.expected {
			t.Errorf("Test %d\n\tExpected %s to equal to %s with input\n\t\t Host: %s; Href:%s\n", i, url, test.expected, test.host, test.href)
		}
	}
}

func TestWithPrefix(t *testing.T) {
	tests := []struct {
		prefix, link string
		expected     bool
	}{
		{"http://wow.com", "http://wow.com/hey", true},
		{"http://wow.com/", "https://hey", false},
	}
	for i, test := range tests {
		f := withPrefix(test.prefix)
		if res := f(test.link); res != test.expected {
			t.Errorf("Test %d\n\tExpected %v to equal to %v with input\n\t\t Prefix: %s; Link:%s\n", i, res, test.expected, test.prefix, test.link)
		}
	}
}
