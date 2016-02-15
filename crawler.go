package html2

import (
	"fmt"
	"net/url"
	"io/ioutil"
	"strings"
	"os"
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
	FilterURL(url, parent *url.URL) bool
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
	if !os.IsPathSeparator(path[len(path)-1]) {
		path = fmt.Sprintf("%v%c",path, os.PathSeparator)
	}
	target := path + d.Content.Filename
	err := ioutil.WriteFile(target, d.Content.bytes.Bytes(), 0644)
	if err != nil {
		fmt.Printf("WHAAAAAAT!? %v\n", err)
	}
}

func (d *DownloadableResult) Process() CrawlResult {
	return &DeadEnd{}
}

type PageResult struct {
	Title string
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

func (c *Crawler) Crawl(urlString string) (result CrawlResult) {
	fmt.Printf("Crawling %v\n", urlString)
	u, err := url.Parse(urlString)
	if err != nil {
		return &Failure{err}
	}
	download, err := c.Fetcher.Fetch(u)
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

		ir := &PageResult{Title:parser.Title}
		for _, link := range parser.Links  {
			for _, rule := range c.Config.Filters.FollowRules {
				if rule.Follow(link.URL) {
					ir.Links = append(ir.Links, link)
				}
			}
		}

		for _, img := range parser.Images {
			realURL, err := combineURL(u, img.URL)
			if err != nil {
				fmt.Printf("Failed to add %v, %v\n", img.URL, err)
			}
			if c.accept(realURL, u) {
				ir.Images = append(ir.Images, &Img{realURL})
			}
		}

		for _, link := range parser.Links {
			realURL, err := combineURL(u, link.URL)
			if err != nil {
				fmt.Printf("Failed to add %v, %v\n", link.URL, err)
			}
			if c.accept(realURL, u) {
				ir.Links = append(ir.Links, &Link{realURL})
			}
		}

		result = ir
	} else {
		fmt.Printf("DOWNLOAD\n")
		result = &DownloadableResult{download}
	}

	return
}


func (c *Crawler) accept(u, parent *url.URL) bool {
	for _, filter := range c.Config.Filters.URLFilters {
		if !filter.FilterURL(u, parent) {
			fmt.Printf("URL filter, %v, failed %v\n", filter, u)
			return false
		}
	}

	return true
}

func combineURL(base, sub *url.URL) (*url.URL, error) {
	if len(sub.Scheme) == 0 {
		if strings.HasPrefix(sub.String(), "/") {
			return url.Parse(base.Host + sub.String())
		} else {
			if strings.HasSuffix(base.String(), "/") {
				return url.Parse(base.String() + sub.String())
			} else {
				pos := strings.LastIndex(base.Path, "/")
				var t string
				if pos < 0 {
					t = base.String()
				} else {
					t = base.String()[:(len(base.String())-1-(len(base.Path)-pos+len(base.RawQuery)))]
				}
				return url.Parse(t + "/" + sub.String())
			}
		}
	} else {
		return sub, nil
	}
}

func getNthIndex(s, t string, n int) int {
	for {
		pos := strings.LastIndex(s, t)
		if pos < 0 {
			return pos
		}
		n--
		if n == 0 {
			return pos
		}

		s = s[:pos]
	}
}
