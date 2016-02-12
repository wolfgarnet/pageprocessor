package html2

import (
	"testing"
	"fmt"
	"reflect"
)

func TestCrawler_1(t *testing.T) {
	config := NewConfiguration()
	crawler := NewCrawler(config)
	result := crawler.Crawl("http://www.ejbyurterne.dk/index.php")

	switch r := result.(type) {
	case *Failure:
		fmt.Printf("FAILED: %v\n", r.Error.Error())

	case *PageResult:
		println("REUSLT I:", len(r.Images))
		println("REUSLT L:", len(r.Links))
	}
}

func TestCrawler_2(t *testing.T) {
	config := NewConfiguration()
	crawler := NewCrawler(config)
	result := crawler.Crawl("http://www.ejbyurterne.dk/graphics/logo.jpg")

	fmt.Printf("TYPE: %v\n", reflect.TypeOf(result))

	switch r := result.(type) {
	case *Failure:
		fmt.Printf("FAILED: %v\n", r.Error.Error())

	case *PageResult:
		r.Process()
		println("REUSLT:", len(r.Images))

	case *DownloadableResult:
		println("REUSLT:", len(r.Content.bytes.Bytes()))
		r.Download("temp/")
	}
}
