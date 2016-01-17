package html2

import (
	"fmt"
	"net/url"
)

type Configuration struct {
	Filters []Filter
}

type Crawler struct {
	Fetcher *Fetcher
	Config *Configuration
}

type Filter interface {
	FilterURL(url *url.URL) bool
	FilterFile(*Download) bool
}

func NewCrawler(config *Configuration) *Crawler {
	return &Crawler{&Fetcher{config}, config}
}

type CrawlResult interface {
	Process() CrawlResult
}

type DownloadableResult struct {
	Content *Download
}

func (d *DownloadableResult) Process() CrawlResult {
	return &DeadEnd{}
}

type PageResult struct {
	Links  []*Link
	Images []*Img
}

func (d *PageResult) Process() CrawlResult {
	return nil
}

type DeadEnd struct {

}

func (d *PageResult) DeadEnd() CrawlResult {
	return nil
}

func (c *Crawler) Crawl(url string) (result CrawlResult) {
	download, err := c.Fetcher.Fetch(url)
	if err != nil {
		fmt.Errorf("Failed, %v\n", err.Error())
		return
	}

	if isDocumentType(download) {
		parser, err := ParseHtml(download.bytes.String())
		if err != nil {
			fmt.Errorf("Failed, %v\n", err.Error())
			return
		}

		result = PageResult{parser.Links, parser.Images}

		fmt.Printf("Links : %v\n", len(parser.Links))
		fmt.Printf("Images: %v\n", len(parser.Images))
	} else {
		result := DownloadableResult{download}
		result.Content = download
	}

	return
}


