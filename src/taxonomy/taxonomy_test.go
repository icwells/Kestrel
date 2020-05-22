// Tests taxonomy methods

package taxonomy

import (
	"strings"
	"testing"
)

func testtaxa(s []string, found bool, nas int) *Taxonomy {
	// Loads data from slice into taxonomy struct
	t := NewTaxonomy()
	t.Kingdom = s[0]
	t.Phylum = s[1]
	t.Class = s[2]
	t.Order = s[3]
	t.Family = s[4]
	t.Genus = s[5]
	t.Species = s[6]
	t.Found = found
	t.Nas = nas
	return t
}

func taxaSlice() []*Taxonomy {
	// Initializes slice of test taxonomies
	var ret []*Taxonomy
	ret = append(ret, testtaxa([]string{"Animalia", "Chordata", "Reptilia", "Squamata", "Anguidae", "Abronia", "Abronia graminea"}, true, 0))
	ret = append(ret, testtaxa([]string{"Animalia", "Chordata", "Reptilia", "Na", "Helodermatidae", "NA", "Heloderma suspectum"}, false, 2))
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

func compareTaxonomies(t *testing.T, e, a *Taxonomy) {
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

func TestCopy(t *testing.T) {
	expected := taxaSlice()
	for _, i := range expected {
		a := NewTaxonomy()
		a.Copy(i)
		compareTaxonomies(t, i, a)
	}
}

func TestCountNAs(t *testing.T) {
	expected := taxaSlice()
	for _, i := range expected {
		e := i.Nas
		i.CountNAs()
		if e != i.Nas {
			t.Errorf("Actual NA count %d does not equal expected: %d", i.Nas, e)
		}
	}
}

func TestCheckLevel(t *testing.T) {
	expected := taxaSlice()
	for _, i := range expected {
		class := i.checkLevel(strings.ToLower(i.Class), false)
		genus := i.checkLevel(strings.ToUpper(i.Genus), false)
		species := i.checkLevel(i.Species, true)
		if class != i.Class {
			t.Errorf("Actual class %s does not equal expected: %s", class, i.Class)
		} else if genus != i.Genus {
			t.Errorf("Actual genus %s does not equal expected: %s", genus, i.Genus)
		} else if species != i.Species {
			t.Errorf("Actual species %s does not equal expected: %s", species, i.Species)
		}
	}
}

func TestSetLevel(t *testing.T) {
	expected := testtaxa([]string{"Animalia", "NA", "Reptilia", "Squamata", "Anguidae", "NA", "Abronia graminea"}, false, 0)
	a := NewTaxonomy()
	taxa := map[string]string{"kingdom": " Animalia", "phylum": "na", "class": "Reptilia", "order": "Squamata", "family": "Anguidae ", "genus": "a", "species": " Abronia graminea "}
	for k, v := range taxa {
		a.SetLevel(k, v)
	}
	compareTaxonomies(t, expected, a)
}
