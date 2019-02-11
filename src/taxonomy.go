// Deifnes taxonomy struct and methods

package main

import (
	"strings"
)

type taxonomy struct {
	kingdom string
	phylum  string
	class   string
	order   string
	family  string
	genus   string
	species string
	source  string
	found   bool
}

func newTaxonomy() taxonomy {
	// Initializes taxonomy struct
	var t taxonomy
	t.kingdom = "NA"
	t.phylum = "NA"
	t.class = "NA"
	t.order = "NA"
	t.family = "NA"
	t.genus = "NA"
	t.species = "NA"
	t.source = "NA"
	t.found = false
	return t
}

func (t *taxonomy) String() string {
	// Returns formatted string without source
	var ret []string
	for _, i := range []string{t.kingdom, t.phylum, t.class, t.order, t.family, t.genus, t.species} {
		ret = append(ret, i)
	}
	return strings.Join(ret, ",")
}
