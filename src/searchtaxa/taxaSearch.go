// Performs taxonomy search for Kestrel program

package searchtaxa

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

func fillLevel(t1, t2 string) string {
	// Returns non-NA value
	if strings.ToUpper(t1) == "NA" && strings.ToUpper(t2) != "NA" {
		t1 = t2
	}
	return t1
}

func fillTaxonomy(t, x taxonomy) taxonomy {
	// Replaces NAs in t with values from x
	t.kingdom = fillLevel(t.kingdom, x.kingdom)
	t.phylum = fillLevel(t.phylum, x.phylum)
	t.class = fillLevel(t.class, x.class)
	t.order = fillLevel(t.order, x.order)
	t.family = fillLevel(t.family, x.family)
	t.genus = fillLevel(t.genus, x.genus)
	t.species = fillLevel(t.species, x.species)
	return t
}

func (s *searcher) setTaxonomy(key, s1, s2 string, t map[string]taxonomy) {
	// Sets taxonomy in searcher map
	if len(s2) > 0 {
		s.terms[key].sources[s2] = t[s2].source
		if t[s1].nas != 0 {
			// Attempt to resolve gaps
			t[s1] = fillTaxonomy(t[s1], t[s2])
		}
	}
	s.terms[key].taxonomy.copyTaxonomy(t[s1])
	s.terms[key].sources[s1] = s.terms[key].taxonomy.source
}

func (s *searcher) getMatch(k string, taxa map[string]taxonomy) bool {
	// Compares results and determines if there has been a match
	ret := false
	var k1, k2 string
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
		} else {
			// Return value with fewest NAs
			min := 8
			for key, v := range taxa {
				if v.nas < min {
					min = v.nas
					k1 = key
				}
			}
		}
	} else if len(taxa) == 1 {
		for key := range taxa {
			k1 = key
		}
	}
	if len(k1) > 0 {
		s.setTaxonomy(k, k1, k2, taxa)
		ret = true
	}
	return ret
}

func checkMatch(taxa map[string]taxonomy, source string, t taxonomy) map[string]taxonomy {
	// Appends t to taxonomy if a match was found
	if t.found == true && t.nas <= 2 {
		taxa[source] = t
	}
	return taxa
}

func (s *searcher) writeResults(mut *sync.RWMutex, k string, found bool) {
	// Manages mutex call and writes to appropriate output file
	mut.Lock()
	if found == true {
		s.writeMatches(k)
	} else {
		// Write missed queries to file
		s.writeMisses(k)
	}
	mut.Unlock()
}

func (s *searcher) searchTerm(wg *sync.WaitGroup, mut *sync.RWMutex, k string) {
	// Performs api search for given term
	defer wg.Done()
	var found bool
	l := strings.Count(s.terms[k].term, "%20") + 1
	for l >= 1 {
		taxa := make(map[string]taxonomy)
		// Search IUCN, NCBI, Wikipedia, and EOL
		taxa = checkMatch(taxa, "IUCN", s.searchIUCN(k))
		taxa = checkMatch(taxa, "NCBI", s.searchNCBI(k))
		taxa = checkMatch(taxa, "EOL", s.searchEOL(k))
		taxa = checkMatch(taxa, "WIKI", s.searchWikipedia(k))
		if len(taxa) >= 1 {
			found = s.getMatch(k, taxa)
		}
		if found == false && l != 1 {
			// Remove first word and try again
			idx := strings.Index(s.terms[k].term, "%20")
			s.terms[k].term = s.terms[k].term[idx+3:]
			l = strings.Count(s.terms[k].term, "%20") + 1
		} else {
			break
		}
	}
	if found == false {
		// Reset term to original
		s.terms[k].term = k
		if s.service.err == nil {
			// Perform selenium search if service is running
			found = s.getSearchResults(k)
		}
	}
	s.writeResults(mut, k, found)
}

func (s *searcher) keySlice() []string {
	// Returns slice of map keys
	var ret []string
	for k := range s.terms {
		ret = append(ret, k)
	}
	return ret
}

func searchTaxonomies(outfile string, searchterms map[string]*Term) {
	// Manages API and selenium searches
	var wg sync.WaitGroup
	var mut sync.RWMutex
	s := newSearcher(false)
	fmt.Println("\n\tSearching for taxonomy matches...")
	if s.service.err == nil {
		defer s.service.stop()
	}
	// Concurrently perform api search
	fmt.Println("\n\tPerforming taxonomy search...")
	for idx, i := range s.keySlice() {
		wg.Add(1)
		go s.searchTerm(&wg, &mut, i)
		fmt.Printf("\tDispatched %d of %d terms.\r", idx+1, len(s.terms))
		if idx%10 == 0 {
			// Pause after 10 to avoid swamping apis
			time.Sleep(time.Second)
		}
		if idx > 1 && idx%200 == 0 {
			// Pause to avoid using all available RAM
			wg.Wait()
		}
	}
	// Wait for remainging processes
	fmt.Println("\n\tWaiting for search results...")
	wg.Wait()
	fmt.Printf("\n\tFound matches for a total of %d queries.\n", s.matches)
	fmt.Printf("\tCould not find matches for %d queries.\n", s.fails)
}
