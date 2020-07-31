// Performs taxonomy search for Kestrel program

package searchtaxa

import (
	"fmt"
	"github.com/icwells/kestrel/src/kestrelutils"
	"github.com/icwells/kestrel/src/taxonomy"
	"github.com/icwells/kestrel/src/terms"
	"os"
	"strings"
	"sync"
	"time"
)

func (s *searcher) setTaxonomy(key, s1, s2 string, t map[string]*taxonomy.Taxonomy) {
	// Sets taxonomy in searcher map
	if len(s2) > 0 {
		if t[s1].Nas != 0 {
			// Attempt to resolve gaps
			s.hier.FillTaxonomy(t[s1])
		}
	}
	s.terms[key].Taxonomy.Copy(t[s1])
}

func (s *searcher) getMatch(k string, taxa map[string]*taxonomy.Taxonomy) bool {
	// Compares results and determines if there has been a match
	ret := false
	var k1, k2, s1, s2 string
	var score int
	if len(taxa) > 1 {
		// Score each pair
		s := newScorer()
		s.setScores(taxa)
		s1, s2, score = s.getMax()
		if len(s1) > 0 {
			// Store key of most complete match and url of supporting match
			if taxa[s1].Nas <= taxa[s2].Nas {
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
				if v.Nas < min {
					min = v.Nas
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
		if score >= 7 || strings.ToLower(s.terms[k].Taxonomy.Species) == strings.ToLower(k) {
			s.terms[k].Confirm()
		}
		ret = true
	}
	return ret
}

func checkMatch(taxa map[string]*taxonomy.Taxonomy, source string, t *taxonomy.Taxonomy) map[string]*taxonomy.Taxonomy {
	// Appends t to taxonomy if a match was found
	if t.Found && t.Nas <= 2 {
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

func (s *searcher) searchCorpus(t *terms.Term) bool {
	// Compares search term to existing taxonomy corpus
	species := t.Term
	if k, ex := s.common[t.Term]; ex {
		species = k
	}
	if match, ex := s.taxa[species]; ex {
		t.Taxonomy.Copy(match)
		t.Confirmed = true
		return true
	}
	return false
}

func (s *searcher) dispatchTerm(k string) bool {
	// Performs api search for given term
	var found bool
	l := strings.Count(s.terms[k].Term, kestrelutils.SPACE) + 1
	for l >= 1 {
		if s.corpus {
			found = s.searchCorpus(s.terms[k])
		}
		if !found {
			taxa := make(map[string]*taxonomy.Taxonomy)
			// Search IUCN, NCBI, Wikipedia, and EOL
			taxa = checkMatch(taxa, "IUCN", s.searchIUCN(k))
			taxa = checkMatch(taxa, "NCBI", s.searchNCBI(k))
			taxa = checkMatch(taxa, "EOL", s.searchEOL(k))
			taxa = checkMatch(taxa, "WIKI", s.searchWikipedia(k))
			if len(taxa) >= 1 {
				found = s.getMatch(k, taxa)
			}
		}
		if !found && s.service.err == nil {
			// Perform selenium search if service is running
			found = s.getSearchResults(k)
		}
		if !found && l != 1 {
			// Remove first word and try again
			idx := strings.Index(s.terms[k].Term, "%20")
			s.terms[k].Term = strings.TrimSpace(s.terms[k].Term[idx+3:])
			l = strings.Count(s.terms[k].Term, "%20") + 1
		} else {
			if l == 1 {
				// Reset term
				s.terms[k].Term = k
			}
			break
		}
	}
	return found
}

func (s *searcher) searchTerm(wg *sync.WaitGroup, mut *sync.RWMutex, k string) {
	// Performs api search for given and corrected term
	defer wg.Done()
	var found bool
	for idx, i := range []string{s.terms[k].Term, s.terms[k].Corrected} {
		if !found && len(i) > 0 {
			if idx == 1 {
				// Set corrected term as term
				s.terms[k].Term, s.terms[k].Corrected = s.terms[k].Corrected, s.terms[k].Term
			}
			found = s.dispatchTerm(k)
			if idx == 1 && !found {
				// Reset original search term
				s.terms[k].Term, s.terms[k].Corrected = s.terms[k].Corrected, s.terms[k].Term
			}
		}
	}
	s.writeResults(mut, k, found)
}

func (s *searcher) searchDone() {
	// Removes previously completed searches
	var completed int
	if s.done.Length() > 0 {
		for k := range s.terms {
			if ex, _ := s.done.InSet(s.terms[k].Term); ex {
				delete(s.terms, k)
				completed++
			}
		}
	}
	if completed > 0 {
		fmt.Printf("\tFound %d terms in previous output.\n", completed)
	}
}

func SearchTaxonomies(outfile string, searchterms map[string]*terms.Term, proc int, nocorpus bool) {
	// Manages API and selenium searches
	var wg sync.WaitGroup
	var mut sync.RWMutex
	count := 1
	s := newSearcher(outfile, searchterms, nocorpus, false)
	if s.service.err == nil {
		defer s.service.stop()
	}
	// Concurrently perform api search
	fmt.Println("\n\tPerforming taxonomy search...")
	s.searchDone()
	if len(s.terms) > 0 {
		fmt.Println("\tPerforming API search...")
		for k := range s.terms {
			wg.Add(1)
			go s.searchTerm(&wg, &mut, k)
			fmt.Printf("\tDispatched %d of %d terms.\r", count, len(s.terms))
			if count%10 == 0 {
				// Pause after 10 to avoid swamping apis
				time.Sleep(time.Second)
			}
			if count%proc == 0 {
				// Pause to avoid using all available RAM
				wg.Wait()
			}
			count++
		}
		// Wait for remainging processes
		fmt.Println("\n\tWaiting for search results...")
		wg.Wait()
	}
	fmt.Printf("\n\tFound matches for a total of %d queries.\n", s.matches)
	fmt.Printf("\tCould not find matches for %d queries.\n", s.fails)
	if s.fails == 0 {
		// Remove unused missed file
		os.Remove(s.missed)
	}
}
