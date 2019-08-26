package sitemap

import "testing"

func TestConstructURL(t *testing.T) {
	tests := []struct {
		host, href, expected string
	}{
		{"http://wow.com", "/hey", "http://wow.com/hey"},
		{"http://wow.com/", "hey", "http://wow.com/hey"},
		{"http://wow.com", "hey", "http://wow.com/hey"},
		{"http://wow.com/", "/hey", "http://wow.com/hey"},
	}
	for i, test := range tests {
		if url := href(test.href, test.host); url != test.expected {
			t.Errorf("Test %d\n\tExpected %s to equal to %s with input\n\t\t Host: %s; Href:%s\n", i, url, test.expected, test.host, test.href)
		}
	}
}
