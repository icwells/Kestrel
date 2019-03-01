// Tests taxaSearch functions

package main

import (
	"strconv"
	"testing"
)

func getTestSearcher() searcher {
	// Returns initialized searcher
	s := newSearcher(true)
	terms := newExtractInput()
	for _, i := range terms {
		if i.status == "" {
			t := newTerm(i.query)
			t.term = i.term
			s.terms[i.term] = &t
		}
	}
	return s
}

func TestFillLevel(t *testing.T) {
	levels := [][]string{
		{"Animalia", "Chordata", "Animalia"},
		{"NA", "Orthoptera", "Orthoptera"},
	}
	for _, i := range levels {
		a := fillLevel(i[0], i[1])
		if a != i[2] {
			t.Errorf("Actual filled value %s does not equal expected: %s", a, i[2])
		}
	}
}

func TestFillTaxonomy(t *testing.T) {
	full := testtaxa([]string{"Animalia", "Chordata", "Reptilia", "Squamata", "Anguidae", "Abronia", "Abronia graminea"}, true, 0)
	taxa := taxaSlice()
	for _, i := range taxa {
		i = fillTaxonomy(i, full)
		i.countNAs()
		if i.nas != 0 {
			t.Errorf("Taxonomy for %s was not filled", i.species)
		}
	}
}

func checkSetTaxonomy(t *testing.T, s searcher, key, k1, k2 string, taxa map[string]taxonomy) {
	compareTaxonomies(t, s.terms[key].taxonomy, taxa[k1])
	if s.terms[key].sources[k1] != taxa[k1].source {
		t.Errorf("Actual source %s does not equal expected: %s", s.terms[key].sources[k1], taxa[k1].source)
	}
	if k2 != "" && s.terms[key].sources[k2] != taxa[k2].source {
		t.Errorf("Actual source %s does not equal expected: %s", s.terms[key].sources[k2], taxa[k2].source)
	}
}

func testTaxonomy() (map[string]taxonomy, []string) {
	// Returns taxonomy map for testing
	var keys []string
	taxa := make(map[string]taxonomy)
	sli := taxaSlice()
	for idx, i := range sli {
		k := strconv.Itoa(idx)
		i.source = k
		taxa = checkMatch(taxa, k, i)
		if _, ex := taxa[k]; ex == true {
			keys = append(keys, k)
		}
	}
	return taxa, keys
}

func TestSetTaxonomy(t *testing.T) {
	s := getTestSearcher()
	taxa, keys := testTaxonomy()
	s.setTaxonomy("Fish", keys[0], keys[1], taxa)
	checkSetTaxonomy(t, s, "Fish", keys[0], keys[1], taxa)
	s.setTaxonomy("Piping Guan", keys[1], keys[0], taxa)
	checkSetTaxonomy(t, s, "Piping Guan", keys[1], "", taxa)
}

func TestGetMatch(t *testing.T) {
	s := getTestSearcher()
	taxa, _ := testTaxonomy()
	a := s.getMatch("Fish", taxa)
	if a != true {
		t.Error("Taxonomy match not found.")
	}
	delete(taxa, "0")
	a = s.getMatch("Fish", taxa)
	if a != true {
		t.Error("Taxonomy match not found.")
	}
}

func TestCheckMatch(t *testing.T) {
	count := 0
	taxa := make(map[string]taxonomy)
	s := taxaSlice()
	for idx, i := range s {
		k := strconv.Itoa(idx)
		taxa := checkMatch(taxa, k, i)
		_, ex := taxa[k]
		if i.found == false || i.nas > 2 {
			if ex == true {
				t.Error("Taxonomy erroneously passed checkMatch.")
			}
		} else {
			if ex == false {
				t.Error("Taxonomy erroneously failed checkMatch.")
			} else {
				count++
			}
		}
	}
	if count != len(taxa) {
		t.Errorf("Count of checkMatch output is incorrect.")
	}
}

func TestKeySlice(t *testing.T) {
	s := getTestSearcher()
	a := s.keySlice()
	if len(a) != len(s.terms) {
		t.Errorf("Actual keySlice length %d does not equal expected: %d", len(a), len(s.terms))
	}
}
