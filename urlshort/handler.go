package urlshort

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		uri := r.URL.Path
		v, ok := pathsToUrls[uri]
		if ok {
			http.Redirect(w, r, v, http.StatusSeeOther)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(handler)
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYAML, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYAML)
	return MapHandler(pathMap, fallback), nil
}

type pathURL struct {
	P string `yaml:"path"`
	U string `yaml:"url"`
}

func buildMap(sl []pathURL) map[string]string {
	m := make(map[string]string)
	for _, v := range sl {
		m[v.P] = v.U
	}
	return m
}

func parseYAML(yml []byte) ([]pathURL, error) {
	var paths []pathURL
	if err := yaml.Unmarshal(yml, &paths); err != nil {
		return nil, err
	}
	return paths, nil
}
