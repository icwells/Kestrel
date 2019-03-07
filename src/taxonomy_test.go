// Tests taxonomy methods

package main

import (
	"strings"
	"testing"
)

func testtaxa(s []string, found bool, nas int) taxonomy {
	// Loads data from slice into taxonomy struct
	t := newTaxonomy()
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

func taxaSlice() []taxonomy {
	// Initializes slice of test taxonomies
	var ret []taxonomy
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
	compareLevels(t, e.kingdom, a.kingdom, "kingdom")
	compareLevels(t, e.phylum, a.phylum, "phylum")
	compareLevels(t, e.class, a.class, "class")
	compareLevels(t, e.order, a.order, "order")
	compareLevels(t, e.family, a.family, "family")
	compareLevels(t, e.genus, a.genus, "genus")
	compareLevels(t, e.species, a.species, "species")
	if a.found != e.found {
		t.Errorf("Actual found value %v does not equal expected: %v", a.found, e.found)
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
		class := i.checkLevel(strings.ToLower(i.class), false)
		genus := i.checkLevel(strings.ToUpper(i.genus), false)
		species := i.checkLevel(strings.Split(i.species, " ")[1], true)
		if class != i.class {
			t.Errorf("Actual class %s does not equal expected: %s", class, i.class)
		} else if genus != i.genus {
			t.Errorf("Actual genus %s does not equal expected: %s", genus, i.genus)
		} else if species != i.species {
			t.Errorf("Actual species %s does not equal expected: %s", species, i.species)
		}
	}
}

func TestCheckTaxa(t *testing.T) {
	expected := taxaSlice()
	for idx, i := range expected {
		a := newTaxonomy()
		a.copyTaxonomy(i)
		if idx == 3 {
			a.kingdom = "Metazoa"
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
