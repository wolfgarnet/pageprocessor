package html2

import "net/url"

type SizeFilter struct {
	Minimum, Maximum int
}

func (sf *SizeFilter) FilterURL(url *url.URL) bool {
	return true
}

func (sf *SizeFilter) FilterFile(download *Download) bool {
	size := len(download.bytes)
	if size < sf.Minimum || size > sf.Maximum {
		return false
	}

	return true
}

type URLFilter struct  {
	Host string

}
