package main

import "github.com/mishuk-sk/gopher/sitemap"

func main() {
	_, err := sitemap.Map("http://miu.by", 5)
	if err != nil {
		panic(err)
	}
}
