// Defines taxonomy struct and methods

package taxonomy

import (
	"github.com/icwells/go-tools/strarray"
	"strings"
)

func fillLevel(t1, t2 string) string {
	// Returns non-NA value
	if strings.ToUpper(t1) == "NA" && strings.ToUpper(t2) != "NA" {
		t1 = t2
	}
	return t1
}

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
	Nas     int
	levels  []string
}

func NewTaxonomy() *Taxonomy {
	// Initializes taxonomy struct
	t := new(Taxonomy)
	t.Kingdom = "NA"
	t.Phylum = "NA"
	t.Class = "NA"
	t.Order = "NA"
	t.Family = "NA"
	t.Genus = "NA"
	t.Species = "NA"
	t.Source = "NA"
	t.Found = false
	t.Nas = 7
	t.levels = []string{"kingdom", "phylum", "class", "order", "family", "genus", "species"}
	return t
}

func (t *Taxonomy) String() string {
	// Returns formatted string without source
	var ret []string
	for _, i := range []string{t.Kingdom, t.Phylum, t.Class, t.Order, t.Family, t.Genus, t.Species} {
		ret = append(ret, i)
	}
	return strings.Join(ret, ",")
}

func (t *Taxonomy) CopyTaxonomy(x *Taxonomy) {
	// Deep copies x to t
	t.Kingdom = x.Kingdom
	t.Phylum = x.Phylum
	t.Class = x.Class
	t.Order = x.Order
	t.Family = x.Family
	t.Genus = x.Genus
	t.Species = x.Species
	t.Source = x.Source
	t.Found = x.Found
	t.Nas = x.Nas
}

func (t *Taxonomy) CountNAs() {
	// Rechecks nas
	nas := 0
	for _, i := range []string{t.Kingdom, t.Phylum, t.Class, t.Order, t.Family, t.Genus, t.Species} {
		if strings.ToUpper(i) == "NA" {
			nas++
		}
	}
	t.Nas = nas
}

func (t *Taxonomy) checkLevel(l string, sp bool) string {
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
				l = t.Genus + " " + strings.ToLower(l)
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

func (t *Taxonomy) CheckTaxa() {
	// Checks formatting
	t.CountNAs()
	if t.Nas <= 2 && strings.ToUpper(t.Genus) != "NA" {
		t.Found = true
		if strings.ToLower(t.Kingdom) == "metazoa" {
			// Correct NCBI kingdom
			t.Kingdom = "Animalia"
		} else {
			t.Kingdom = t.checkLevel(t.Kingdom, false)
		}
		t.Phylum = t.checkLevel(t.Phylum, false)
		t.Class = t.checkLevel(t.Class, false)
		t.Order = t.checkLevel(t.Order, false)
		t.Family = t.checkLevel(t.Family, false)
		t.Genus = t.checkLevel(t.Genus, false)
		t.Species = t.checkLevel(t.Species, true)
	}
}

func (t *Taxonomy) SetLevel(key, value string) {
	// Sets level denoted by key with value
	value = strings.TrimSpace(value)
	if strings.Contains(value, "[") == false && strings.ToUpper(value) != "NA" && len(value) > 1 {
		switch key {
		case "kingdom":
			t.Kingdom = value
		case "phylum":
			t.Phylum = value
		case "class":
			t.Class = value
		case "order":
			t.Order = value
		case "family":
			t.Family = value
		case "genus":
			t.Genus = value
		case "species":
			t.Species = value
		}
	}
}

func (t *Taxonomy) IsLevel(s string) string {
	// Returns formatted string if s is a taxonomic level
	s = strings.TrimSpace(strings.ToLower(strings.Replace(s, ":", "", -1)))
	for _, i := range t.levels {
		if i == s {
			return s
		}
	}
	return ""
}

func (t *Taxonomy) FillTaxonomy(x *Taxonomy) {
	// Replaces NAs in t with values from x
	t.Kingdom = fillLevel(t.Kingdom, x.Kingdom)
	t.Phylum = fillLevel(t.Phylum, x.Phylum)
	t.Class = fillLevel(t.Class, x.Class)
	t.Order = fillLevel(t.Order, x.Order)
	t.Family = fillLevel(t.Family, x.Family)
	t.Genus = fillLevel(t.Genus, x.Genus)
	t.Species = fillLevel(t.Species, x.Species)
}
