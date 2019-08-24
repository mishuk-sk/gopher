package linkparser

import (
	"fmt"
	"io"
	"strings"
	"unicode"

	"golang.org/x/net/html"
)

//Link represents both <a>-tag's href and text inside
type Link struct {
	Href string
	Text string
}

//Parse returns slice of Link that were read from r
func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("Error parsing html via x/net/html. Err - %s", err)
	}
	var links []Link
	aTags := getNodes(doc)
	for _, n := range aTags {
		links = append(links, readATag(n))
	}
	return links, nil
}

func getNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}
	var nodes []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, getNodes(c)...)
	}
	return nodes
}

func readATag(n *html.Node) Link {
	var link Link
	text := strings.Builder{}
	for _, a := range n.Attr {
		if a.Key == "href" {
			link.Href = a.Val
		}
	}
	getText(n, &text)
	split := strings.FieldsFunc(
		text.String(),
		func(r rune) bool { return unicode.IsSpace(r) },
	)
	link.Text = strings.Join(split, " ")
	return link
}

func getText(n *html.Node, buf *strings.Builder) {
	if n.Type == html.TextNode {
		buf.WriteString(n.Data)
	} else if n.Type != html.ElementNode {
		return
	} else {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			getText(c, buf)
		}
	}
}
