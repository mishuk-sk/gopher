package sitemap

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/mishuk-sk/gopher/linkparser"
)

// Map returns site map as xml for provided url
func Map(url string, depth int) ([]byte, error) {
	pages, err := pagesMap(url, depth)
	if err != nil {
		return nil, err
	}
	urls := make([]string, 0, len(pages))
	for k := range pages {
		urls = append(urls, k)
	}
	sort.Strings(urls)
	for _, u := range urls {
		fmt.Println(u)
	}
	return nil, nil
}

func pagesMap(url string, depth int) (map[string]struct{}, error) {
	seen := make(map[string]struct{})
	var q map[string]struct{}
	nq := map[string]struct{}{
		url: struct{}{},
	}
	host, err := getHost(url)
	if err != nil {
		return nil, err
	}
	for i := 0; i < depth; i++ {
		if len(nq) == 0 {
			break
		}
		q, nq = nq, make(map[string]struct{})

		for u := range q {
			if _, ok := seen[u]; ok {
				continue
			}
			seen[u] = struct{}{}
			links := get(u, host)
			for _, link := range links {
				nq[link] = struct{}{}
			}
		}

	}
	return seen, nil
}

func getHost(u string) (string, error) {
	r, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return "", fmt.Errorf("Error performing request to url %s. Err - %s", u, err)
	}
	baseURL := &url.URL{
		Scheme: r.URL.Scheme,
		Host:   r.URL.Host,
	}
	return baseURL.String(), nil
}

func get(u, base string) []string {
	resp, err := http.Get(u)
	if err != nil {
		log.Printf("Error getting page from %s. Err - %s\n", u, err)
		return []string{}
	}
	defer resp.Body.Close()
	links, err := constructLinks(resp.Body, base)
	if err != nil {
		log.Printf("Error getting links for page on %s. Err - %s\n", u, err)
		return []string{}
	}
	return filter(links, withPrefix(base))
}

func constructLinks(r io.Reader, base string) ([]string, error) {
	links, err := linkparser.Parse(r)
	if err != nil {
		return nil, err
	}
	var resp []string
	for _, link := range links {
		resp = append(resp, href(link.Href, base))
	}
	return resp, nil
}

func href(l, base string) string {
	switch {
	case strings.HasPrefix(l, "/"):
		return base + l
	case strings.HasPrefix(l, "http"):
		return l
	default:
		return ""
	}
}

func filter(links []string, keepFn func(string) bool) []string {
	var resp []string
	for _, link := range links {
		if keepFn(link) {
			resp = append(resp, link)
		}
	}
	return resp
}

func withPrefix(prefix string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, prefix)
	}
}
