// Defines taxonomy struct and methods

package taxonomy

import (
	"github.com/icwells/go-tools/strarray"
	"strings"
	"unicode"
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
	t.Source = ""
	t.Found = false
	t.Nas = 7
	t.levels = []string{"kingdom", "phylum", "class", "order", "family", "genus", "species"}
	return t
}

func (t *Taxonomy) Slice(id, db string) []string {
	// Returns slice of taxonomy
	var ret []string
	if id != "" {
		ret = append(ret, id)
	}
	for _, i := range []string{t.Kingdom, t.Phylum, t.Class, t.Order, t.Family, t.Genus, t.Species, t.Source} {
		ret = append(ret, strings.Replace(i, `"`, "", -1))
	}
	if db != "" {
		ret = append(ret, db)
	}
	return ret
}

func (t *Taxonomy) String() string {
	// Returns formatted string
	return strings.Join(t.Slice("", ""), ",")
}

func (t *Taxonomy) Copy(x *Taxonomy) {
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

func (t *Taxonomy) SpeciesCaps(name string) string {
	// Properly capitalizes species name
	name = strings.TrimSpace(strings.ToLower(name))
	s := strings.Split(name, " ")
	if len(s) > 1 {
		// Save with genus capitalized and species in lower case
		var builder strings.Builder
		builder.WriteString(strarray.TitleCase(s[0]))
		for _, i := range s[1:] {
			builder.WriteByte(' ')
			builder.WriteString(i)
		}
		return builder.String()
	} else {
		return strarray.TitleCase(name)
	}
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

func (t *Taxonomy) removePunctuation(s string) string {
	// Removes punctuation from line
	var ret strings.Builder
	for _, i := range []rune(s) {
		if !unicode.IsPunct(i) {
			ret.WriteRune(i)
		}
	}
	return ret.String()
}

func (t *Taxonomy) checkLevel(l string, sp bool) string {
	// Returns formatted name
	if strings.ToUpper(l) != "NA" {
		l = strings.Replace(l, ",", "", -1)
		if sp == false {
			if strings.Contains(l, " ") {
				l = strings.Split(l, " ")[0]
			}
			l = strarray.TitleCase(l)
		} else {
			// Get binomial with proper capitalization
			if strings.Contains(l, ".") {
				// Remove genus abbreviations
				l = strings.TrimSpace(l[strings.Index(l, ".")+1:])
			}
			if !strings.Contains(l, " ") {
				l = t.Genus + " " + strings.ToLower(l)
			} else {
				s := strings.Split(l, " ")
				l = strarray.TitleCase(s[0]) + " " + strings.ToLower(s[1])
			}
		}
	} else {
		// Standardize NAs
		l = strings.ToUpper(l)
	}
	return t.removePunctuation(l)
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
	if len(value) > 1 && !strings.Contains(value, "[") && strings.ToUpper(value) != "NA" {
		key = strings.ToLower(key)
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

func (t *Taxonomy) latinToEnglish(key string) string {
	// Translates Latin level names to English
	switch key {
	case "regnum":
		key = "kingdom"
	// Phylum, genus, and specie are the same
	case "classis":
		key = "class"
	case "ordo":
		key = "order"
	case "familia":
		key = "family"
	}
	return key
}

func (t *Taxonomy) IsLevel(s string, translate bool) string {
	// Returns formatted string if s is a taxonomic level
	s = strings.TrimSpace(strings.ToLower(strings.Replace(s, ":", "", -1)))
	if translate {
		s = t.latinToEnglish(s)
	}
	if strarray.InSliceStr(t.levels, s) {
		return s
	}
	return ""
}

func (t *Taxonomy) ContainsLevel(s string) string {
	// Returns level if s contains one
	s = strings.ToLower(s)
	for _, i := range t.levels {
		if strings.Contains(s, i) {
			return i
		}
	}
	return ""
}
