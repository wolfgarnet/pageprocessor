package html2

import "fmt"

type Crawler struct {
	Fetcher *Fetcher
}

func NewCrawler() *Crawler {
	return &Crawler{&Fetcher{}}
}

func (c *Crawler) Crawl(url string) {
	d, err := c.Fetcher.Fetch(url)
	if err != nil {
		fmt.Errorf("Failed, %v\n", err.Error())
		return
	}

	d.Display()
	fmt.Printf("DOWNLOAD: %+v\n", d)

	if isDocumentType(d) {
		parser, err := ParseHtml(d.bytes.String())
		if err != nil {
			fmt.Errorf("Failed, %v\n", err.Error())
			return
		}

		fmt.Printf("Links: %v\n", len(parser.Links))

		for i, l := range parser.Links {
			fmt.Printf("LINK [%v] %v\n", i, l.Link)
		}

		for i, img := range parser.Images {
			fmt.Printf("IMAGE [%v] %v\n", i, img.Url)
		}
	}
}


