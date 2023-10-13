package telegraph

import (
	"bytes"
	"io"
	"strings"

	"golang.org/x/net/html"
)

type FilterFunc func(domNode *html.Node) bool

// ContentFormat transforms data to a DOM-based format to represent the content of the page.
func ContentFormat(data any, filters ...FilterFunc) (n []Node, err error) {
	var dst *html.Node

	switch src := data.(type) {
	case string:
		dst, err = html.Parse(strings.NewReader(src))
	case []byte:
		dst, err = html.Parse(bytes.NewReader(src))
	case io.Reader:
		dst, err = html.Parse(src)
	default:
		return nil, ErrInvalidDataType
	}

	if err != nil {
		return nil, err
	}

	if node := domToNode(dst.FirstChild, filters...); node != nil {
		n = append(n, node)
	}

	return n, nil
}

func domToNode(domNode *html.Node, filters ...FilterFunc) any {
	for _, filter := range filters {
		if filter(domNode) {
			return nil
		}
	}
	if domNode.Type == html.TextNode {
		return domNode.Data
	}

	if domNode.Type != html.ElementNode {
		return nil
	}

	nodeElement := new(NodeElement)

	switch strings.ToLower(domNode.Data) {
	case "a", "aside", "b", "blockquote", "br", "code", "em", "figcaption", "figure", "h3", "h4", "hr", "i",
		"iframe", "img", "li", "ol", "p", "pre", "s", "strong", "u", "ul", "video":
		nodeElement.Tag = domNode.Data

		for i := range domNode.Attr {
			switch strings.ToLower(domNode.Attr[i].Key) {
			case "href", "src":
				nodeElement.Attrs = map[string]string{domNode.Attr[i].Key: domNode.Attr[i].Val}
			default:
				continue
			}
		}
	}

	for child := domNode.FirstChild; child != nil; child = child.NextSibling {
		if node := domToNode(child, filters...); node != nil {
			nodeElement.Children = append(nodeElement.Children, domToNode(child, filters...))
		}
	}

	return nodeElement
}
