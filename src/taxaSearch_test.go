// Tests taxaSearch functions

package main

import (
	"testing"
)

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
