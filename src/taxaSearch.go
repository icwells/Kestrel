// Performs taxonomy search for Kestrel program

package main

import (
	"fmt"
	//"os"
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

func (s *searcher) getMatch(k string, last int, taxa map[string]taxonomy) bool {
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
		} else if last == 1 {
			// Return value with fewest NAs
			min := 8
			for key, v := range taxa {
				if v.nas < min {
					min = v.nas
					k1 = key
				}
			}
		}
	} else if last == 1 {
		// Only accept single match for last search
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

func (s *searcher) searchTerm(wg *sync.WaitGroup, mut *sync.RWMutex, k string) {
	// Performs api search for given term
	defer wg.Done()
	var found bool
	l := strings.Count(s.terms[k].term, "%20") + 1
	for l >= 1 {
		taxa := make(map[string]taxonomy)
		// Search IUCN, NCBI, Wikipedia, and EOL
		taxa = checkMatch(taxa, "IUCN", s.searchIUCN(k))
		//taxa = checkMatch(taxa, "NCBI", s.searchNCBI(k))
		//taxa = checkMatch(taxa, "EOL", s.searchEOL(k))
		//taxa = checkMatch(taxa, "WIKI", s.searchWikipedia(k))
		if len(taxa) >= 1 {
			found = s.getMatch(k, l, taxa)
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
	if found == true {
		mut.Lock()
		s.writeMatches(k)
		mut.Unlock()
	} else {
		// Record missed keys and reset term to original
		s.terms[k].term = k
		s.misses = append(s.misses, k)
	}
}

func (s *searcher) keySlice() []string {
	// Returns slice of map keys
	var ret []string
	for k := range s.terms {
		ret = append(ret, k)
	}
	return ret
}

func searchTaxonomies(start time.Time) {
	// Manages API and selenium searches
	var wg sync.WaitGroup
	var mut sync.RWMutex
	s := newSearcher(false)
	s.termMap(*infile)
	// Concurrently perform api search
	fmt.Println("\n\tPerforming API based taxonomy search...")
	for idx, i := range s.keySlice() {
		//s.misses = append(s.misses, k)
		wg.Add(1)
		go s.searchTerm(&wg, &mut, i)
		if idx%10 == 0 {
			// Pause after 10 to avoid swamping apis
			time.Sleep(2 * time.Second)
		}
		fmt.Printf("\tDispatched %d of %d terms.\r", idx+1, len(s.terms))
	}
	// Wait for remainging processes
	fmt.Println("\n\tWaiting for search results...")
	wg.Wait()
	fmt.Printf("\tFound matches for %d queries.\n", s.matches)
	fmt.Printf("\tCurrent run time: %v\n\n", time.Since(start))
	/*if len(s.misses) > 0 {
		// Perform selenium search on misses
		f := s.matches
		service, browser, err := getBrowser(*firefox)
		if err == nil {
			defer service.Stop()
			defer browser.Quit()
			for idx, i := range s.misses {
				res := s.seleniumSearch(browser, i)
				// Parse search results concurrently
				wg.Add(1)
				go s.getSearchResults(&wg, &mut, res, i)
				if idx%10 == 0 {
					time.Sleep(2 * time.Second)
				}
				fmt.Printf("\tDispatched %d of %d missed terms.\r", idx+1, len(s.misses))
			}
			wg.Wait()
			fmt.Println("\n\tWaiting for search results...")
			fmt.Printf("\tFound matches for %d missed queries.\n\n", s.matches-f)
		} else {
			fmt.Printf("\t[Error] Could not initialize Selenium server: %v\n", err)
			fmt.Println("\n\tWriting misses to file...")
			for _, i := range s.misses {
				s.writeMisses(i)
			}
		}
	}*/
	fmt.Printf("\tFound matches for a total of %d queries.\n", s.matches)
	fmt.Printf("\tCould not find matches for %d queries.\n\n", s.fails)
}
