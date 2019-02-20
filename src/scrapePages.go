// Contains web scraping functions

package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func getPage(url string) (io.Reader, bool) {
	// Wraps http request, returns io.Reader and true if successful
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
	ret.scrapeWiki(url)
	return ret
}

func (s *searcher) esearch(term string) string {
	// Returns taxonomy ID for search term
	var id string
	url := fmt.Sprintf("%sesearch.fcgi?db=Taxonomy&term=%s&api_key=%s", s.urls.ncbi, term, s.keys["NCBI"])
	page, err := goquery.NewDocument(url)
	if err == nil {
		q := page.Find("Id")
		if len(q.Text()) >= 1 {
			id = q.Text()
		}
	}
	return id
}

func (s *searcher) espell(term string) string {
	// Checks spelling of term
	term = strings.Replace(term, " ", "%20", -1)
	url := fmt.Sprintf("%sespell.fcgi?db=Taxonomy&term=%s&api_key=%s", s.urls.ncbi, term, s.keys["NCBI"])
	page, err := goquery.NewDocument(url)
	if err == nil {
		q := page.Find("correctedquery")
		if _, err := strconv.Atoi(q.Text()); err == nil {
			term = strings.Replace(q.Text(), " ", "%20", -1)
		}
	}
	return term
}

func (s *searcher) searchNCBI(k string) taxonomy {
	// Searches NCBI for species ID and uses id to query taxonomy
	ret := newTaxonomy()
	res := s.espell(s.terms[k].term)
	if len(res) > 0 {
		id := s.esearch(res)
		if len(id) > 0 {
			url := fmt.Sprintf("%sefetch.fcgi?db=Taxonomy&id=%s$retmode=xml&api_key=%s", s.urls.ncbi, id, s.keys["NCBI"])
			ret.scrapeNCBI(url)
		}
	}
	return ret
}

func (s *searcher) searchIUCN(k string) taxonomy {
	// Seaches IUCN Red List for match
	ret := newTaxonomy()
	if _, ex := s.keys["IUCN"]; ex == true {
		url := fmt.Sprintf("%s%s?token=%s", s.urls.iucn, s.terms[k].term, s.keys["IUCN"])
		result, pass := getPage(url)
		if pass == true {
			ret.scrapeIUCN(result, url)
		}
	}
	fmt.Println(ret.String())
	return ret
}
