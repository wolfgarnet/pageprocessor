package html2

import "testing"

func TestCrawler(t *testing.T) {
	config := &Configuration{}
	crawler := NewCrawler(config)
	result := crawler.Crawl("http://www.ejbyurterne.dk/index.php")
	println("REUSLT:", result.Images)
}
