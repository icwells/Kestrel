// Performs taxonomy search for Kestrel program

package searchtaxa

import (
	"fmt"
	"github.com/icwells/kestrel/src/kestrelutils"
	"github.com/icwells/kestrel/src/taxonomy"
	"github.com/icwells/kestrel/src/terms"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

func (s *searcher) setTaxonomy(k, key string, t map[string]*taxonomy.Taxonomy) {
	// Sets taxonomy in searcher map
	if t[key].Nas != 0 {
		// Attempt to resolve gaps
		s.hier.FillTaxonomy(t[key])
	}
	s.terms[k].Taxonomy.Copy(t[key])
}

func (s *searcher) getMatch(k string, taxa map[string]*taxonomy.Taxonomy) bool {
	// Compares results and determines if there has been a match
	ret := false
	var key, s1, s2 string
	var score int
	if len(taxa) > 1 {
		// Score each pair
		s := newScorer()
		s.setScores(taxa)
		s1, s2, score = s.getMax()
		if len(s1) > 0 {
			// Store key of most complete match and url of supporting match
			if taxa[s1].Nas <= taxa[s2].Nas {
				key = s1
			} else {
				key = s2
			}
		} else {
			// Return value with fewest NAs
			min := 8
			for name, v := range taxa {
				if v.Nas < min {
					min = v.Nas
					key = name
				}
			}
		}
	} else if len(taxa) == 1 {
		for name := range taxa {
			key = name
		}
	}
	if len(key) > 0 {
		s.setTaxonomy(k, key, taxa)
		if score >= 7 || strings.ToLower(s.terms[k].Taxonomy.Species) == strings.ToLower(k) {
			s.terms[k].Confirm()
		} else if v := s.corpusMatch(k); v != "" {
			s.terms[k].Confirm()
		}
		ret = true
	}
	return ret
}

func checkMatch(taxa map[string]*taxonomy.Taxonomy, t *taxonomy.Taxonomy) map[string]*taxonomy.Taxonomy {
	// Appends t to taxonomy if a match was found
	if t.Found && t.Nas <= 2 {
		taxa[t.Source] = t
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

func (s *searcher) corpusMatch(name string) string {
	// Returns key for taxa map if found
	if k, ex := s.common[name]; ex {
		name = k
	}
	if _, ex := s.taxa[name]; ex {
		return name
	}
	return ""
}

func (s *searcher) searchCorpus(t *terms.Term) bool {
	// Compares search term to existing taxonomy corpus
	if k := s.corpusMatch(t.Term); k != "" {
		t.Taxonomy.Copy(s.taxa[k])
		t.Confirmed = true
		return true
	} else {
		// Attempt to find fuzzy match
		matches := fuzzy.RankFindFold(t.Term, s.names)
		if matches.Len() > 0 {
			sort.Sort(matches)
			if matches[0].Distance <= int(float64(len(t.Term))*0.1) {
				if k := s.corpusMatch(matches[0].Target); k != "" {
					t.Taxonomy.Copy(s.taxa[k])
					return true
				}
			}
		}
	}
	return false
}

func (s *searcher) wordCount(k string) int {
	// Returns number of words
	return strings.Count(s.terms[k].Term, kestrelutils.SPACE) + 1
}

func (s *searcher) dispatchTerm(k string) bool {
	// Performs api search for given term
	var found bool
	for !found {
		l := s.wordCount(k)
		if s.corpus {
			found = s.searchCorpus(s.terms[k])
		}
		if !found {
			taxa := make(map[string]*taxonomy.Taxonomy)
			// Search IUCN, NCBI, Wikipedia, and EOL
			taxa = checkMatch(taxa, s.searchIUCN(k))
			taxa = checkMatch(taxa, s.searchNCBI(k))
			taxa = checkMatch(taxa, s.searchEOL(k))
			taxa = checkMatch(taxa, s.searchWikipedia(k))
			taxa = checkMatch(taxa, s.searchWikiSpecies(k))
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
			l = s.wordCount(k)
		} else if l == 1 {
			// Reset term
			s.terms[k].Term = k
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
