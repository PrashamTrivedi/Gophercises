package link

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

//Link Represents a link in HTML
type Link struct {
	Href string
	Text string
}

//Parse will take an HTML Document and will return a slice of links from that HTML
// or an error
func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	nodes := linkNodes(doc)
	var links []Link
	for _, node := range nodes {
		links = append(links, buildLink(node))
	}

	return links, nil
}

func buildLink(n *html.Node) Link {
	var link Link
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			link.Href = attr.Val
		}
	}
	link.Text = text(n)
	return link
}

func text(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}
	var textData string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		textData += text(c)
	}
	return strings.Join(strings.Fields(textData), " ")
}

func linkNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}
	var ret []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, linkNodes(c)...)
	}
	return ret
}
