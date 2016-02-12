package html2

import (
	"fmt"
	"net/url"
	"io/ioutil"
)

type Filters struct  {
	URLFilters []URLFilter
	FileFilters []FileFilter
	FollowRules []FollowRule
}

func NewFilters() *Filters {
	filters := &Filters{
		URLFilters:make([]URLFilter, 0),
		FileFilters:make([]FileFilter, 0),
		FollowRules:make([]FollowRule, 0),
	}

	return filters
}

type Configuration struct {
	Filters *Filters
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Filters:NewFilters(),
	}
}

type Crawler struct {
	Fetcher *Fetcher
	Config *Configuration
}

type URLFilter interface {
	FilterURL(url *url.URL) bool
}

type FileFilter interface {
	FilterFile(*Download) bool
}

type FollowRule interface {
	Follow(url *url.URL) bool
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

func (d *DownloadableResult) Download(path string) {
	target := path + d.Content.Filename
	fmt.Printf("Writing to %v\n", target)
	err := ioutil.WriteFile(target, d.Content.bytes.Bytes(), 0644)
	if err != nil {
		fmt.Printf("WHAAAAAAT!? %v\n", err)
	}
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

		ir := &PageResult{}
		for _, link := range parser.Links  {
			for _, rule := range c.Config.Filters.FollowRules {
				if rule.Follow(link.URL) {
					ir.Links = append(ir.Links, link)
				}
			}
		}

		for _, img := range parser.Images {
			ir.Images = append(ir.Images, img)
		}

		for _, link := range parser.Links {
			ir.Links = append(ir.Links, link)
		}

		result = ir
	} else {
		fmt.Printf("DOWNLOAD\n")
		result = &DownloadableResult{download}
	}

	return
}


