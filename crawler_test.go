package html2

import (
	"testing"
	"fmt"
	"reflect"
	"net/url"
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

func TestCombineURL(t *testing.T) {
	tests := []struct {
		url string
		anchor string
		expected string
	}{
		{"http://www.a.dk/index.php?pid=main", "index.php?pid=next", "http://www.a.dk/index.php?pid=next"},
	}

	for i, test := range tests {
		l1, err := url.Parse(test.url)
		if err != nil {
			t.Errorf("Failed test #%v, %v", i, err)
		}
		l2, err := url.Parse(test.anchor)
		if err != nil {
			t.Errorf("Failed test #%v, %v", i, err)
		}
		u, err := combineURL(l1, l2)
		if err != nil {
			t.Errorf("Failed test #%v, %v", i, err)
		}

		if test.expected != u.String() {
			t.Errorf("Failed test #%v,\nExpected: %v\nActual  : %v", i, test.expected, u.String())
		}

		fmt.Printf("%v == %v\n", test.expected, u.String())
	}

}
