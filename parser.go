package html2

import (
	"golang.org/x/net/html"
	"io"
	"strings"
	"net/url"
	"fmt"
)


type EntityType int

/*
const (
	Link EntityType = iota + 1
	Image
	Title
)
*/

type Entity interface {
	Children() []Entity
}

type (
	Img struct {
		Url string
	}

	Link struct {
		URL *url.URL
	}
)

type HtmlParser struct {
	Links  []*Link
	Images []*Img
}

func ParseHtml(html string) (*HtmlParser, error) {
	parser := &HtmlParser{}
	return parser, parser.parseString(html)
}

func (hp *HtmlParser) parseString(src string) error {
	return hp.parseReader(strings.NewReader(src))
}

func (hp *HtmlParser) parseReader(src io.Reader) error {
	doc, err := html.Parse(src)
	if err != nil {
		return err
	}

	hp.parse(doc)

	return nil
}

func (hp *HtmlParser) parse(node *html.Node) {
	if node.Type == html.ElementNode {
		switch node.Data {
		case "a":
			println("HEJ", node.Attr[0].Val)
			link := hp.parseLink(node)

			hp.Links = append(hp.Links, link)
		case "img":
			img := hp.parseImg(node)
			hp.Images = append(hp.Images, img)
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		hp.parse(c)
	}
}

func findAttr(key string, attributes []html.Attribute) string {
	for _, a := range attributes {
		if a.Key == key {
			return a.Val
		}
	}

	return ""
}

func (hp *HtmlParser) parseLink(node *html.Node) *Link {
	href := findAttr("href", node.Attr)
	fmt.Printf("HREF= %v\n", href)
	if href == "" {
		return nil
	}

	url, err := url.Parse(href)
	if err != nil {
		fmt.Printf("%v skipped, %err\n", href, err)
		return nil
	}

	return &Link{url}
}

func (hp *HtmlParser) parseImg(node *html.Node) *Img {
	src := findAttr("src", node.Attr)

	if src == "" {
		return nil
	}

	return &Img{src}
}
