// Contains web scraping functions

package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"net/http"
	"strconv"
	"strings"
)

func getPage(url string) ([]byte, bool) {
	// Wraps http request, returns io.Reader and true if successful
	var ret []byte
	pass := false
	resp, err := http.Get(url)
	if err == nil {
		// Convert reader to byte slice
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(resp.Body)
		if err == nil {
			ret = buf.Bytes()
			pass = true
		}
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
	if _, ex := s.keys["NCBI"]; ex == true {
		res := s.espell(s.terms[k].term)
		if len(res) > 0 {
			id := s.esearch(res)
			if len(id) > 0 {
				url := fmt.Sprintf("%sefetch.fcgi?db=Taxonomy&id=%s$retmode=xml&api_key=%s", s.urls.ncbi, id, s.keys["NCBI"])
				ret.scrapeNCBI(url)
			}
		}
	}
	return ret
}

func (s *searcher) getHID(tid string) string {
	// Returns hierarchy id from EOL
	var ret string
	url := fmt.Sprintf("%s%sxml?id=%s&vetted=1&key=%s", s.urls.eol, s.urls.pages, tid, s.keys["EOL"])
	page, err := goquery.NewDocument(url)
	if err == nil {
		page.Find("taxonConcept").EachWithBreak(func(i int, r *goquery.Selection) bool {
			if r.Find("taxonRank").Text() == "species" {
				// Skip incomplete taxonomies
				hid := r.Find("identifier").Text()
				if _, er := strconv.Atoi(hid); er == nil {
					ret = hid
					return false
				}
			}
			return true
		})
	}
	return ret
}

func (s *searcher) getTID(term string) string {
	// Gets taxon id from EOL search api
	var ret string
	score := len(term)
	query := percentDecode(term)
	url := fmt.Sprintf("%s%sxml?q=%s&vetted=1&key=%s", s.urls.eol, s.urls.search, term, s.keys["EOL"])
	page, err := goquery.NewDocument(url)
	if err == nil {
		page.Find("result").EachWithBreak(func(i int, r *goquery.Selection) bool {
			// Iterate though all results
			id := r.Find("id").Text()
			if _, err := strconv.Atoi(id); err == nil {
				// Examine all valid ids
				title := r.Find("title").Text()
				if fuzzy.MatchFold(query, title) == true {
					// Keep scientific name match
					ret = id
					return false
				}
				content := strings.Split(r.Find("content").Text(), ";")
				for _, i := range content {
					// Examine each content entry seperately
					dist := fuzzy.LevenshteinDistance(query, i)
					if dist == 0 {
						// Keep perfect match
						ret = id
						return false
					} else if dist < score {
						// Store best match
						score = dist
						ret = id
					}
				}
			}
			return true
		})
	}
	return ret
}

func (s *searcher) searchEOL(k string) taxonomy {
	// Searches EOL for taxon id, hierarchy entry id, and taxonomy
	ret := newTaxonomy()
	if _, ex := s.keys["EOL"]; ex == true {
		tid := s.getTID(k)
		if len(tid) >= 1 {
			hid := s.getHID(tid)
			if len(hid) >= 1 {
				// Switch to json for easier scraping of larger results
				url := fmt.Sprintf("%s%sjson?id=%s&vetted=1&key=%s", s.urls.eol, s.urls.hier, hid, s.keys["EOL"])
				result, pass := getPage(url)
				if pass == true {
					ret.scrapeEOL(result, url)
				}
			}
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
	return ret
}
