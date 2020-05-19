// Defines taxonomy struct and methods

package kestrelutils

import (
	"github.com/icwells/go-tools/strarray"
	"strings"
)

type Taxonomy struct {
	Kingdom string
	Phylum  string
	Class   string
	Order   string
	Family  string
	Genus   string
	Species string
	Source  string
	Found   bool
	nas     int
	levels  []string
}

func NewTaxonomy() *Taxonomy {
	// Initializes taxonomy struct
	t := make(Taxonomy)
	t.Kingdom = "NA"
	t.Phylum = "NA"
	t.Class = "NA"
	t.Order = "NA"
	t.Family = "NA"
	t.Genus = "NA"
	t.Species = "NA"
	t.Source = "NA"
	t.Found = false
	t.nas = 7
	t.levels = []string{"kingdom", "phylum", "class", "order", "family", "genus", "species"}
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

func (t *taxonomy) CopyTaxonomy(x taxonomy) {
	// Deep copies x to t
	t.kingdom = x.kingdom
	t.phylum = x.phylum
	t.class = x.class
	t.order = x.order
	t.family = x.family
	t.genus = x.genus
	t.species = x.species
	t.source = x.source
	t.found = x.found
	t.nas = x.nas
}

func (t *taxonomy) CountNAs() {
	// Rechecks nas
	nas := 0
	for _, i := range []string{t.Kingdom, t.Phylum, t.Class, t.Order, t.Family, t.Genus, t.Species} {
		if strings.ToUpper(i) == "NA" {
			nas++
		}
	}
	t.nas = nas
}

func (t *taxonomy) checkLevel(l string, sp bool) string {
	// Returns formatted name
	if strings.ToUpper(l) != "NA" {
		l = strings.Replace(l, ",", "", -1)
		if sp == false {
			if strings.Contains(l, " ") == true {
				l = strings.Split(l, " ")[0]
			}
			l = strarray.TitleCase(l)
		} else {
			// Get binomial with proper capitalization
			if strings.Contains(l, ".") == true {
				// Remove genus abbreviations
				l = strings.TrimSpace(l[strings.Index(l, ".")+1:])
			}
			if strings.Contains(l, " ") == false {
				l = t.genus + " " + strings.ToLower(l)
			} else {
				s := strings.Split(l, " ")
				l = strings.Title(s[0]) + " " + strings.ToLower(s[1])
			}
		}
	} else {
		// Standardize NAs
		l = strings.ToUpper(l)
	}
	return l
}

func (t *taxonomy) CheckTaxa() {
	// Checks formatting
	t.countNAs()
	if t.nas <= 2 && strings.ToUpper(t.genus) != "NA" {
		t.found = true
		if strings.ToLower(t.kingdom) == "metazoa" {
			// Correct NCBI kingdom
			t.kingdom = "Animalia"
		} else {
			t.kingdom = t.checkLevel(t.kingdom, false)
		}
		t.phylum = t.checkLevel(t.phylum, false)
		t.class = t.checkLevel(t.class, false)
		t.order = t.checkLevel(t.order, false)
		t.family = t.checkLevel(t.family, false)
		t.genus = t.checkLevel(t.genus, false)
		t.species = t.checkLevel(t.species, true)
	}
}

func (t *taxonomy) SetLevel(key, value string) {
	// Sets level denoted by key with value
	value = strings.TrimSpace(value)
	if strings.Contains(value, "[") == false && strings.ToUpper(value) != "NA" && len(value) > 1 {
		switch key {
		case "kingdom":
			t.kingdom = value
		case "phylum":
			t.phylum = value
		case "class":
			t.class = value
		case "order":
			t.order = value
		case "family":
			t.family = value
		case "genus":
			t.genus = value
		case "species":
			t.species = value
		}
	}
}

func (t *taxonomy) IsLevel(s string) string {
	// Returns formatted string if s is a taxonomic level
	s = strings.TrimSpace(strings.ToLower(strings.Replace(s, ":", "", -1)))
	for _, i := range t.levels {
		if i == s {
			return s
		}
	}
	return ""
}
