// Contains web scraping functions

package main

import (
	//"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func getPage(url string) (io.Reader, bool) {
	// Wraps http request, returns byte slice and true if successful
	var ret io.Reader
	pass := false
	resp, err := http.Get(url)
	if err == nil {
		ret = resp.Body
		pass = true
	}
	defer resp.Body.Close()
	return ret, pass
}

func (s *searcher) searchWikipedia(k string) taxonomy {
	// Scrapes taxonomy from Wikipedia entry
	ret := newTaxonomy()
	t := s.terms[k].term
	if strings.Contains(t, "%20") == true {
		t = strings.Replace(t, "%20", "_", -1)
	}
	url := s.urls.wiki + t
	result, pass := getPage(url)
	if pass == true {
		ret.scrapeWiki(result, url)
	}
	return ret
}

func (s *searcher) searchIUCN(k string) taxonomy {
	// Seaches IUCN Red List for match
	ret := newTaxonomy()
	key, ex := s.keys["IUCN"]
	if ex == true {
		url := fmt.Sprintf("%s%s?token=%s", s.urls.iucn, s.terms[k].term, s.keys["IUCN"])
		result, pass := getPage(url)
		if pass == true {
			ret.scrapeIUCN(result, url)
		}
	}
	return ret
}
