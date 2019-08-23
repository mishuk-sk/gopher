package linkparser

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

//Link represents both <a>-tag's href and text inside
type Link struct {
	Href string
	Text string
}

//Links returns slice of Link that were read from r
func Links(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("Error parsing html via x/net/html. Err - %s", err)
	}
	var links []Link
	var dfs func(n *html.Node)
	dfs = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			links = append(links, readATag(n))
		} else {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				dfs(c)
			}
		}
	}
	dfs(doc)
	return links, nil
}

func readATag(n *html.Node) Link {
	var link Link
	text := strings.Builder{}
	for _, a := range n.Attr {
		if a.Key == "href" {
			link.Href = a.Val
		}
	}
	var dfs func(n *html.Node)
	dfs = func(n *html.Node) {
		if n.Type == html.TextNode {
			text.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			dfs(c)
		}
	}
	dfs(n)
	link.Text = text.String()
	return link
}
