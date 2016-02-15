package html2

import (
	"net/url"
	"strings"
)

func findFilters(f interface{}, filters *Filters) {

	f1, ok := f.(URLFilter)
	if ok {
		filters.URLFilters = append(filters.URLFilters, f1)
	}

	f2, ok := f.(FileFilter)
	if ok {
		filters.FileFilters = append(filters.FileFilters, f2)
	}

	f3, ok := f.(FollowRule)
	if ok {
		filters.FollowRules = append(filters.FollowRules, f3)
	}
}

// SIZE FILTER

type SizeFilter struct {
	Minimum, Maximum int
}

func (sf *SizeFilter) FilterFile(download *Download) bool {
	if download.bytes.Len() < sf.Minimum || download.bytes.Len() > sf.Maximum {
		return false
	}

	return true
}

// PAGE EXTENSION FILTER

type PageExtensionFilter struct  {
	Extensions []string
	Allowed bool
}

func (pef *PageExtensionFilter) FilterURL(url, parent *url.URL) bool {
	for _, ext := range pef.Extensions {
		if url.Fragment == ext {
			return pef.Allowed
		}
	}

	return !pef.Allowed
}

// KEYWORD RULE

type KeywordRuleFilter struct {
	Whitelist []string
	Blacklist []string
}

func (kw *KeywordRuleFilter) FilterURL(url, parent *url.URL) bool {
	upper := strings.ToUpper(url.String())
	counted := 0
	for _, w := range kw.Whitelist {
		if strings.ToUpper(w) == upper {
			return true
		}

		counted++
	}

	if counted > 0 {
		return false
	}

	for _, b := range kw.Blacklist {
		if strings.ToUpper(b) == upper {
			return false
		}
	}

	return true
}

// No cross site crawl

type NoCrossSiteCrawl struct {
}

func (ncsc *NoCrossSiteCrawl) FilterURL(url, parent *url.URL) bool {
	if url.Host != parent.Host {
		return false
	}

	return true
}
