// Tests taxonomy methods

package kestrelutils

import (
	"strings"
	"testing"
)

func testtaxa(s []string, found bool, nas int) taxonomy {
	// Loads data from slice into taxonomy struct
	t := NewTaxonomy()
	t.kingdom = s[0]
	t.phylum = s[1]
	t.class = s[2]
	t.order = s[3]
	t.family = s[4]
	t.genus = s[5]
	t.species = s[6]
	t.found = found
	t.nas = nas
	return t
}

func taxaSlice() []*Taxonomy {
	// Initializes slice of test taxonomies
	var ret []*Taxonomy
	ret = append(ret, testtaxa([]string{"Animalia", "Chordata", "Reptilia", "Squamata", "Anguidae", "Abronia", "Abronia graminea"}, true, 0))
	ret = append(ret, testtaxa([]string{"Animalia", "Chordata", "Actinopteri", "Na", "Pomacentridae", "NA", "NA saxatilis"}, false, 2))
	ret = append(ret, testtaxa([]string{"Animalia", "Na", "NA", "Orthoptera", "NA", "Acheta", "Acheta domesticus"}, false, 3))
	ret = append(ret, testtaxa([]string{"Animalia", "Chordata", "Aves", "Passeriformes", "Sturnidae", "Acridotheres", "Acridotheres tristis"}, true, 0))
	return ret
}

func compareLevels(t *testing.T, e, a, l string) {
	// Compares values at given level
	if e != a && e != "Na" {
		t.Errorf("Actual %s %s does not equal expected: %s", l, a, e)
	}
}

func compareTaxonomies(t *testing.T, e, a taxonomy) {
	// Compares expected to actual taxonomies
	compareLevels(t, e.Kingdom, a.Kingdom, "kingdom")
	compareLevels(t, e.Phylum, a.Phylum, "phylum")
	compareLevels(t, e.Class, a.Class, "class")
	compareLevels(t, e.Order, a.Order, "order")
	compareLevels(t, e.Family, a.Family, "family")
	compareLevels(t, e.Genus, a.Genus, "genus")
	compareLevels(t, e.Species, a.Species, "species")
	if a.Found != e.Found {
		t.Errorf("Actual found value %v does not equal expected: %v", a.Found, e.Found)
	}
}

func TestCopyTaxonomy(t *testing.T) {
	expected := taxaSlice()
	for _, i := range expected {
		a := newTaxonomy()
		a.copyTaxonomy(i)
		compareTaxonomies(t, i, a)
	}
}

func TestCountNAs(t *testing.T) {
	expected := taxaSlice()
	for _, i := range expected {
		e := i.nas
		i.countNAs()
		if e != i.nas {
			t.Errorf("Actual NA count %d does not equal expected: %d", i.nas, e)
		}
	}
}

func TestCheckLevel(t *testing.T) {
	expected := taxaSlice()
	for _, i := range expected {
		class := i.checkLevel(strings.ToLower(i.Class), false)
		genus := i.checkLevel(strings.ToUpper(i.Genus), false)
		species := i.checkLevel(strings.Split(i.Species, " ")[1], true)
		if class != i.Class {
			t.Errorf("Actual class %s does not equal expected: %s", class, i.Class)
		} else if genus != i.Genus {
			t.Errorf("Actual genus %s does not equal expected: %s", genus, i.Genus)
		} else if species != i.Species {
			t.Errorf("Actual species %s does not equal expected: %s", species, i.Species)
		}
	}
}

func TestCheckTaxa(t *testing.T) {
	expected := taxaSlice()
	for idx, i := range expected {
		a := newTaxonomy()
		a.copyTaxonomy(i)
		if idx == 3 {
			a.Kingdom = "Metazoa"
		}
		a.checkTaxa()
		compareTaxonomies(t, i, a)
	}
}

func TestSetLevel(t *testing.T) {
	expected := testtaxa([]string{"Animalia", "NA", "Reptilia", "Squamata", "Anguidae", "NA", "Abronia graminea"}, false, 0)
	a := newTaxonomy()
	taxa := map[string]string{"kingdom": " Animalia", "phylum": "na", "class": "Reptilia", "order": "Squamata", "family": "Anguidae ", "genus": "a", "species": " Abronia graminea "}
	for k, v := range taxa {
		a.setLevel(k, v)
	}
	compareTaxonomies(t, expected, a)
}

func TestIsLevel(t *testing.T) {
	taxa := newTaxonomy()
	expected := []struct {
		input, expected string
	}{
		{" Species", "species"},
		{"six", ""},
		{"GENUS ", "genus"},
		{"Famil", ""},
	}
	for _, i := range expected {
		a := taxa.isLevel(i.input)
		if a != i.expected {
			t.Errorf("Actual isLevel output %s, does not equal expected: %s", a, i.expected)
		}
	}
}
