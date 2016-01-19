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

func (d *DeadEnd) Process() CrawlResult {
	return nil
}

type Failure struct {
	Error error
}

func (d *Failure) Process() CrawlResult {
	return nil
}

func (c *Crawler) Crawl(url string) (result CrawlResult) {
	download, err := c.Fetcher.Fetch(url)
	if err != nil {
		return &Failure{err}
	}

	if isDocumentType(download) {
		fmt.Printf("PAGE\n")
		parser, err := ParseHtml(download.bytes.String())
		if err != nil {
			fmt.Errorf("Failed, %v\n", err.Error())
			return &Failure{err}
		}

		result = &PageResult{parser.Links, parser.Images}
	} else {
		fmt.Printf("DOWNLOAD\n")
		result = &DownloadableResult{download}
	}

	return
}


