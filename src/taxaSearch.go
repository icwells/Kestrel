// Performs taxonomy search for Kestrel program

package main

import (
	"fmt"
	"strings"
)

func (s *searcher) setTaxonomy(key, s1, s2, source string, t taxonomy) {
	// Sets taxonomy in searcher map
	s.terms[key].taxonomy = t
	s.terms[key].sources[s1] = t.source
	if len(source) > 0 {
		s.terms[key].sources[s2] = source
	}
}

func (s *searcher) getMatch(k string, last bool, taxa map[string]taxonomy) bool {
	// Compares results and determines if there has been a match
	ret := false
	var k1, k2, source string
	if len(taxa) > 1 {
		// Score each pair
		s := newScorer()
		s.setScores(taxa)
		s1, s2 := s.getMax()
		if len(s1) > 0 {
			// Store key of most complete match and url of supporting match
			if taxa[s1].nas <= taxa[s2].nas {
				k1 = s1
				k2 = s2
			} else {
				k1 = s2
				k2 = s1
			}
			source = taxa[k2].source
		}
	} else if last == true {
		// Only accept single match for last search
		for k, v := range taxa {
			k1 = k
		}
	}
	if len(k1) > 0 {
		s.setTaxonomy(k, k1, k2, source, taxa[k1])
		ret = true
	}
	return ret
}

func checkMatch(taxa map[string]taxonomy, source string, t taxonomy) map[string]taxonomy {
	// Appends t to taxonomy if a match was found
	if t.found == true {
		taxa[source] = t
	}
	return taxa
}

func (s *searcher) searchTerm(k string) {
	// Performs api search for given term
	var e taxonomy
	var found, last bool
	searchterm := s.terms[k].term
	l := len(strings.Split(searchterm, " "))
	for l >= 1 {
		taxa := make(map[string]taxonomy)
		if l == 1 {
			last = true
		}
		// Search IUCN, NCBI, Wikipedia, and EOL
		taxa = checkMatch(taxa, "IUCN", s.searchIUCN(k))
		taxa = checkMatch(taxa, "NCBI", s.searchNCBI(k))
		taxa = checkMatch(taxa, "WIKI", s.searchWikipedia(k))
		if len(taxa) < 2 {
			// Prioritize against EOL since their results are not returned in order of relevance
			taxa = checkMatch(taxa, "EOL", s.searchEOL(k))
		}
		if len(taxa) >= 1 {
			found = s.getMatch(k, last, taxa)
		}
		if found == false && last == false {
			// Remove first word and try again
			idx := strings.Index(s.terms[k].term, "%20")
			s.terms[k].term = s.terms[k].term[idx+1:]
			l = len(s.terms[k].term)
		} else {
			break
		}
	}
	if found == true {
		s.writeMatches(k)
	} else {
		// Record missed keys
		s.misses = append(s.misses, k)
	}
}

func searchTaxonomies() {
	// Manages API and selenium searches
	s := newSearcher()
	s.termMap(*infile)
	l := len(s.terms)
	// Concurrently perform api search
	fmt.Println("\n\tPerforming API based taxonomy search...")
	for k := range s.terms {
		go s.searchTerm(k)
		fmt.Printf("\tFound %d of %d matches.\r", s.matches, l)
	}
	// Perform selenium search on misses

}
