// Defines scorer struct and methods for scoring taxonomy comparisons

package searchtaxa

import (
	"github.com/icwells/kestrel/src/taxonomy"
)

type scorer struct {
	scores map[string]map[string]int
}

func newScorer() scorer {
	// Initializes map
	var s scorer
	s.scores = make(map[string]map[string]int)
	return s
}

func (s *scorer) getMax() (string, string) {
	// Returns keys of highest scoring match
	var r1, r2 string
	max := -8
	for key, value := range s.scores {
		for k, v := range value {
			if v > max {
				// Store keys of highest score
				max = v
				r1 = key
				r2 = k
			}
		}
	}
	if max <= 5 {
		// Reject low scoring match
		r1 = ""
		r2 = ""
	}
	return r1, r2
}

func (s *scorer) scoreLevel(t1, t2 string) int {
	// +1 for match, -1 for mismatch, +0 for NA
	if t1 == "NA" || t2 == "NA" {
		return 0
	} else if t1 == t2 {
		return 1
	} else {
		return -1
	}
}

func (s *scorer) score(t1, t2 *taxonomy.Taxonomy) int {
	// Scores each taxonomy
	ret := 0
	ret += s.scoreLevel(t1.Kingdom, t2.Kingdom)
	ret += s.scoreLevel(t1.Phylum, t2.Phylum)
	ret += s.scoreLevel(t1.Class, t2.Class)
	ret += s.scoreLevel(t1.Order, t2.Order)
	ret += s.scoreLevel(t1.Family, t2.Family)
	ret += s.scoreLevel(t1.Genus, t2.Genus)
	ret += s.scoreLevel(t1.Species, t2.Species)
	return ret
}

func (s *scorer) setScores(taxa map[string]*taxonomy.Taxonomy) {
	// Calculate scores for each pairing
	var sources []string
	var t []*taxonomy.Taxonomy
	for key, value := range taxa {
		// Get linked slices to use indeces
		sources = append(sources, key)
		t = append(t, value)
	}
	for start := 0; start < len(sources)-1; start++ {
		k1 := sources[start]
		if _, ex := s.scores[k1]; ex == false {
			s.scores[k1] = make(map[string]int)
		}
		for idx := start + 1; idx < len(sources); idx++ {
			// Compare to each successive taxonomy
			k2 := sources[idx]
			s.scores[k1][k2] = s.score(t[start], t[idx])
		}
	}
}
