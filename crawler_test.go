package html2

import (
	"testing"
	"fmt"
)

func TestCrawler(t *testing.T) {
	config := &Configuration{}
	crawler := NewCrawler(config)
	result := crawler.Crawl("http://www.ejbyurterne.dk/index.php")

	switch r := result.(type) {
	case *Failure:
		fmt.Printf("FAILED: %v\n", r.Error.Error())

	case *PageResult:
		println("REUSLT:", len(r.Images))
	}
}
