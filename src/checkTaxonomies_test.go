// Test checkTaxonomies and hiererchy functions

package main

import (
	"testing"
)

type testcheck struct {
	term		string
	taxonomy	[]string
	match		bool
}

func newtestcheck(t string, s []string, match bool) testcheck {
	// Returns filled struct
	var ret testcheck
	ret.term = t
	ret.taxonomy = s
	ret.match = match
	return ret
}

func getTestTaxa() []testcheck {
	// Returns struct slice of test taxonmies
	var ret []testcheck
	ret = append(ret, newtestcheck("Death Adder", []string{"Animalia", "Chordata", "Reptilia", "NA", "Hydrophiidae", "Acanthophis", "Acanthophis antarcticus"}, true))
	ret = append(ret, newtestcheck("Coopers hawk", []string{"Animalia", "Chordata", "NA", "Accipitriformes", "Accipitridae", "NA", "Accipiter cooperii"}, true))
	ret = append(ret, newtestcheck("Finish goshawk", []string{"Animalia", "Chordata", "Aves", "Falconiformes", "Accipitridae", "Accipiter", "Accipiter gentilis"}, true))
	ret = append(ret, newtestcheck("Giant Fruit Bat", []string{"Animalia", "Chordata", "NA", "Chiroptera", "Pteropodidae", "Acerodon", "Acerodon lucifer"}, true))
	ret = append(ret, newtestcheck("Cairo Spiny Mouse", []string{"Animalia", "Chordata", "Mammalia", "Rodentia", "Muridae", "Acomys", "Acomys cahirinus"}, false))
	return ret
}

func getCurated() curated {
	// Returns test struct of currated taxonomies
	c := newCurated()
	c.taxa["Death adder"] = []string{"Animalia", "Chordata", "Reptilia", "Squamata", "Hydrophiidae", "Acanthophis", "Acanthophis antarcticus"}
	c.taxa["Cooper's hawk"] = []string{"Animalia", "Chordata", "Aves", "Accipitriformes", "Accipitridae", "Accipiter", "Accipiter cooperii"}
	c.taxa["Finnish goshawk"] = []string{"Animalia", "Chordata", "Aves", "Falconiformes", "Accipitridae", "Accipiter", "Accipiter gentilis"}
	c.taxa["Giant fruit bat"] = []string{"Animalia", "Chordata", "Mammalia", "Chiroptera", "Pteropodidae", "Acerodon", "Acerodon lucifer"}
	c.setKeys()
	c.set = true
	return c
}

func TestGetTaxonomy(t *testing.T) {
	c := getCurated()
	cases := getTestTaxa()
	for _, i := range cases {
		_, pass := c.getTaxonomy(i.term)
		if pass != i.match {
			t.Errorf("Actual match value %v does not equal expected: %v", pass, i.match)
		}
	}
}

func getHierarchy(levels []string) hierarchy {
	// Returns filled hiearchy struct
	ret := newHierarchy()
	ret.header = make(map[string]int)
	c := getCurated()
	for _, v := range c.taxa {
		for idx, i := range v {
			if idx > 1 {
				ret.setParent(levels[idx], v[idx-1], i)
			}
		}
	}
	for idx, i := range levels {
		ret.header[i] = idx
	}
	return ret
}

func TestCheckHierarchy(t *testing.T) {
	levels := []string{"Kingdom", "Phylum", "Class", "Order", "Family", "Genus", "Species"}
	h := getHierarchy(levels)
	cases := getTestTaxa()
	for _, c := range cases {
		s := h.checkHierarchy(c.taxonomy)
		for idx, i := range s {
			if i == "NA" {
				t.Errorf("NA found in %s of %s.", levels[idx], c.term)
			}
		}
	}
}
