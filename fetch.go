package html2

import (
	"net/url"
	"strings"
	"net/http"
	"bytes"
	"fmt"
	"mime"
	"path/filepath"
)

type ContentType int

const (
	Binary ContentType = iota
	Html
)

func (c ContentType) String() string {
	switch c {
	case Binary:
		return "Binary"

	case Html:
		return "Html"

	default:
		return "WAAT"
	}
}

type Download struct {
	ContentType ContentType
	Filename string
	Type string
	bytes *bytes.Buffer
}

func (d *Download) Display() {
	fmt.Printf("%v: %v\n", d.Filename, d.Type)
}

func getContentType(header *http.Header) ContentType {
	contentType := header.Get("Content-Type")
	switch contentType {
	case "text/html":
		return Html

	default:
		return Binary
	}
}

type Fetcher struct {
	config *Configuration
}

func isDocumentType(download *Download) bool {
	switch download.ContentType {
	case Html:
		return true

	default:
		return false
	}
}

func (f *Fetcher) Fetch(fileURL *url.URL) (*Download, error) {
	path := fileURL.Path

	segments := strings.Split(path, "/")

	check := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	resp, err := check.Get(fileURL.String()) // add a filter to check redirect

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()


	//disposition := resp.Header.Get("Content-Disposition")
	fmt.Printf("DISPO: %v\n", resp.Header)

	//resp.

	download := &Download{getContentType(&resp.Header), segments[len(segments) - 1], "", &bytes.Buffer{}}
	download.bytes.ReadFrom(resp.Body)
	fmt.Printf("----->%v\n", mime.TypeByExtension(filepath.Ext(download.Filename)))
	download.Type = mime.TypeByExtension(filepath.Ext(download.Filename))

	for _, filter := range f.config.Filters.FileFilters {
		if !filter.FilterFile(download) {
			return nil, fmt.Errorf("File filter, %v, failed")
		}
	}

	println(resp.Status)
	fmt.Printf("---->%v\n", download)

	return download, nil
}
