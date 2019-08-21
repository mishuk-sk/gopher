package urlshort_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mishuk-sk/gopher/urlshort"
)

type testCase struct {
	input      string
	outputPath string
	redirect   bool
}

func TestMapHandler(t *testing.T) {
	tests := []testCase{
		{
			input:      "/http",
			outputPath: "https://godoc.org/net/http/httptest",
			redirect:   true,
		},
		{
			input:      "/https",
			outputPath: "/https",
			redirect:   false,
		},
	}
	paths := buildMap(tests)
	handle := urlshort.MapHandler(paths, defaultMux())
	for _, test := range tests {
		req := httptest.NewRequest("GET", test.input, nil)
		w := httptest.NewRecorder()
		handle(w, req)
		url, err := w.Result().Location()
		if err != nil {
			if test.redirect {
				t.Errorf("Path %s was not redirected  to %s\n", test.input, test.outputPath)
			}
			continue
		}
		path := url.Scheme + "://" + url.Host + url.Path
		if path != test.outputPath && url.Path != test.outputPath {
			t.Errorf("Expected path %s to equal %s\n", path, test.outputPath)
		}
	}
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world!")
	})
	return mux
}

func buildMap(tests []testCase) map[string]string {
	m := make(map[string]string)
	for _, test := range tests {
		if test.redirect {
			m[test.input] = test.outputPath
		}
	}
	return m
}
