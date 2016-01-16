package html2

import (
	"testing"
	"io"
	"os"
	"fmt"
)

func TestFetcher(t *testing.T) {
	url := "http://seventyeight.org/78.png"

	fetcher := &Fetcher{}
	download, err := fetcher.Fetch(url)

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
