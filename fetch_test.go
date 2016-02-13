package html2

import (
	"testing"
	"io"
	"os"
	"fmt"
	"net/url"
)

func TestFetcher(t *testing.T) {
	url1 := "http://seventyeight.org/78.png"
	link, _ := url.Parse(url1)

	fetcher := &Fetcher{}
	download, err := fetcher.Fetch(link)

	if err != nil {
		fmt.Errorf(err.Error())
	}

	filename := "YEAH.png"

	file, err := os.Create(filename)

	if err != nil {
		fmt.Errorf(err.Error())
	}
	defer file.Close()

	io.Copy(file, download.bytes)
}
