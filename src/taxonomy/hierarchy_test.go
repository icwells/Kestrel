// Tests hierarchy struct

package taxonomy

import (
	"testing"
)

func hierSlice() []*Taxonomy {
	// Initializes slice of test taxonomies
	var ret []*Taxonomy
	ret = append(ret, testtaxa([]string{"Animalia", "Chordata", "Reptilia", "Squamata", "Anguidae", "Abronia", "Abronia graminea"}, true, 0))
	ret = append(ret, testtaxa([]string{"Animalia", "Chordata", "Reptilia", "Squamata", "Helodermatidae", "Heloderma", "Heloderma suspectum"}, true, 0))
	ret = append(ret, testtaxa([]string{"Animalia", "Chordata", "Insecta", "Orthoptera", "Gryllidae", "Acheta", "Acheta domesticus"}, true, 0))
	ret = append(ret, testtaxa([]string{"Animalia", "Chordata", "Aves", "Passeriformes", "Sturnidae", "Acridotheres", "Acridotheres tristis"}, true, 0))
	return ret
}

func TestHierarchy(t *testing.T) {
	taxa := hierSlice()
	h := NewHierarchy(taxa)
	for _, i := range taxa {
		if _, ex := h.species[i.Species]; !ex {
			t.Errorf("%s not found in species map.", i.Species)
		} else if v := h.species[i.Species]; v != i.Genus {
			t.Errorf("%s genus %s is incorrect.", i.Species, v)
		} else if _, ex := h.genus[i.Genus]; !ex {
			t.Errorf("%s not found in genus map.", i.Genus)
		} else if v := h.genus[i.Genus]; v != i.Family {
			t.Errorf("%s family %s does not equal %s.", i.Genus, v, i.Family)
		} else if _, ex := h.family[i.Family]; !ex {
			t.Errorf("%s not found in family map.", i.Family)
		} else if v := h.family[i.Family]; v != i.Order {
			t.Errorf("%s order %s does not equal %s.", i.Family, v, i.Order)
		} else if _, ex := h.order[i.Order]; !ex {
			t.Errorf("%s not found in order map.", i.Order)
		} else if v := h.order[i.Order]; v != i.Class {
			t.Errorf("%s class %s does not equal %s.", i.Order, v, i.Class)
		} else if _, ex := h.class[i.Class]; !ex {
			t.Errorf("%s not found in class map.", i.Class)
		} else if v := h.class[i.Class]; v != i.Phylum {
			t.Errorf("%s phylum %s does not equal %s.", i.Class, v, i.Phylum)
		} else if _, ex := h.phylum[i.Phylum]; !ex {
			t.Errorf("%s not found in phylum map.", i.Phylum)
		} else if v := h.phylum[i.Phylum]; v != i.Kingdom {
			t.Errorf("%s kingdom %s does not equal %s.", i.Phylum, v, i.Kingdom)
		}
	}
}

func TestFillTaxonomy(t *testing.T) {
	taxa := hierSlice()
	h := NewHierarchy(taxa)
	for _, i := range taxaSlice() {
		h.FillTaxonomy(i)
		if i.Nas != 0 {
			t.Errorf("%s contains %d NAs.", i.Species, i.Nas)
		}
	}
}
