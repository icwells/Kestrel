// Contains web scraping functions

package main

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

func getPage(source, term, key string) {
	// Wraps http request
	

}

func (s *searcher) searchIUCN(k string) taxonomy {
	// Seaches IUCN Red List for match
	ret := newTaxonomy()
	key, ex := s.keys["IUCN"]
	if ex == true {
		result, url = getPage(s.urls.iucn, s.terms[k].term, key)

	return ret
}
