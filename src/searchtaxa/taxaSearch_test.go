// Tests taxaSearch functions

package searchtaxa

import (
	"github.com/icwells/dbIO"
	"github.com/icwells/kestrel/src/taxonomy"
	"github.com/icwells/kestrel/src/terms"
	"strconv"
	"testing"
)

func testtaxa(s []string) *taxonomy.Taxonomy {
	// Loads data from slice into taxonomy struct
	t := taxonomy.NewTaxonomy()
	t.Kingdom = s[0]
	t.Phylum = s[1]
	t.Class = s[2]
	t.Order = s[3]
	t.Family = s[4]
	t.Genus = s[5]
	t.Species = s[6]
	return t
}

func taxaSlice() []*taxonomy.Taxonomy {
	// Initializes slice of test taxonomies
	var ret []*taxonomy.Taxonomy
	ret = append(ret, testtaxa([]string{"Animalia", "Chordata", "Reptilia", "Squamata", "Anguidae", "Abronia", "Abronia graminea"}))
	ret = append(ret, testtaxa([]string{"Animalia", "Chordata", "Reptilia", "Squamata", "Helodermatidae", "Heloderma", "Heloderma suspectum"}))
	ret = append(ret, testtaxa([]string{"Animalia", "Chordata", "Insecta", "Orthoptera", "Gryllidae", "Acheta", "Acheta domesticus"}))
	ret = append(ret, testtaxa([]string{"Animalia", "Chordata", "Aves", "Passeriformes", "Sturnidae", "Acridotheres", "Acridotheres tristis"}))
	return ret
}

func getTestSearcher() searcher {
	// Returns initialized searcher
	var db *dbIO.DBIO
	exp := make(map[string]*terms.Term)
	queries := [][]string{[]string{"abronia Graminea", "Abronia graminea"},
		[]string{"GILA MONSTER", "Gila monster"},
		[]string{"cricket", "Cricket"},
		[]string{"Acridotheres tristis", "Acridotheres tristis"},
	}
	for _, i := range queries {
		t := terms.NewTerm(i[0])
		t.Term = i[1]
		exp[i[1]] = t
	}
	s := newSearcher(db, "", exp, true, true)
	return s
}

/*func TestGetMatch(t *testing.T) {
	s := getTestSearcher()
	taxa, _ := testTaxonomy()
	a := s.getMatch("cricket", taxa)
	if a != true {
		t.Error("Taxonomy match not found.")
	}
	delete(taxa, "0")
	a = s.getMatch("Fish", taxa)
	if a != true {
		t.Error("Taxonomy match not found.")
	}
}*/

func TestCheckMatch(t *testing.T) {
	count := 0
	taxa := make(map[string]*taxonomy.Taxonomy)
	s := taxaSlice()
	for idx, i := range s {
		k := strconv.Itoa(idx)
		taxa := checkMatch(taxa, i)
		_, ex := taxa[k]
		if i.Found == false || i.Nas > 2 {
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
