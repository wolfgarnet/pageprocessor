package html2

import (
	"testing"
	"fmt"
)

func TestCrawler(t *testing.T) {
	config := &Configuration{}
	crawler := NewCrawler(config)
	result := crawler.Crawl("http://www.ejbyurterne.dl/index.php")

	switch r := result.(type) {
	case *Failure:
		fmt.Printf("FAILED: %v\n", r.Error.Error())
	}

	_, ok := result.(*PageResult)
	if !ok {
		t.Errorf("Failed. Not page result")
	}
//	println("REUSLT:", len(r.Images))
}
