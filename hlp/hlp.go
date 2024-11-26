package hlp

import (
	"io"
	"strings"

	"github.com/Junior-Green/gophercises/queue"
	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func extractLink(node *html.Node) Link {
	var text, href string

	for _, attr := range node.Attr {
		if attr.Key == "href" {
			href = attr.Val
			break
		}
	}

	traverse(node, func(node *html.Node) {
		if node.Type == html.TextNode {
			text += node.Data + " "
		}
	})

	return Link{
		href,
		strings.TrimSpace(text),
	}
}

func traverse(root *html.Node, do func(*html.Node)) {
	var q queue.Queue[*html.Node]
	q.Push(root)

	for q.Length() > 0 {
		n := q.Dequeue()

		do(n)

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			q.Push(child)
		}
	}

}

func ExtractLinks(file io.Reader) ([]Link, error) {
	root, err := html.Parse(file)

	if err != nil {
		return nil, err
	}

	links := make([]Link, 0)

	traverse(root, func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			links = append(links, extractLink(n))
		}
	})

	return links, nil
}
