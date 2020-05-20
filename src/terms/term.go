// Defines term struct

package terms

import (
	"github.com/icwells/go-tools/strarray"
	"github.com/icwells/kestrel/src/kestrelutils"
	"github.com/icwells/kestrel/src/taxonomy"
	"regexp"
	"strings"
	"unicode"
)

type Term struct {
	Queries  []string
	Term     string
	Status   string
	Taxonomy *taxonomy.Taxonomy
	Sources  map[string]string
}

func newTerm(query string) *Term {
	// Returns initialized term
	t := new(Term)
	if len(query) > 0 {
		t.addQuery(query)
	}
	t.Taxonomy = taxonomy.NewTaxonomy()
	t.Sources = make(map[string]string)
	return t
}

func (t *Term) String() string {
	// Returns formatted string
	var ret []string
	ret = append(ret, kestrelutils.PercentDecode(t.Term))
	ret = append(ret, t.Taxonomy.String())
	// Append url or NA for each source
	for _, i := range []string{"IUCN", "NCBI", "WIKI", "EOL", "ITIS"} {
		s, ex := t.Sources[i]
		if ex == true {
			ret = append(ret, s)
		} else {
			ret = append(ret, "NA")
		}
	}
	return strings.Join(ret, ",")
}

func (t *Term) addQuery(query string) {
	// Appends to query slice
	t.Queries = append(t.Queries, query)
}

func (t *Term) checkRunes() {
	// Removes puntuation and numbers from term
	var name strings.Builder
	for _, i := range []rune(t.Term) {
		if unicode.IsLetter(i) == true || unicode.IsSpace(i) == true {
			// Remove punctuation and numbers
			name.WriteRune(i)
		} else if i == '.' || i == '-' || i == '\'' {
			name.WriteRune(i)
		}
	}
	t.Term = name.String()
	// Double check starting and ending runes for escaped punctuation
	if len(t.Term) > 0 {
		if t.Term[0] == '.' || t.Term[0] == '-' {
			t.Term = t.Term[1:]
		}
		if t.Term[len(t.Term)-1] == '-' {
			t.Term = t.Term[:len(t.Term)-1]
		}
	}
}

func (t *Term) sliceTerm(p1, p2 string) {
	// Removes item from between 2 puntuation marks
	idx := strings.Index(t.Term, p1)
	ind := strings.LastIndex(t.Term, p2)
	if idx >= 0 && idx < ind {
		// Drop item between punctuation
		if ind == len(t.Term)-1 {
			t.Term = t.Term[:idx]
		} else if idx == 0 {
			t.Term = t.Term[ind+1:]
		} else {
			t.Term = t.Term[:idx] + t.Term[ind+1:]
		}
		t.Term = strings.TrimSpace(t.Term)
	} else {
		// Remove puntuation
		t.Term = strings.Replace(t.Term, p1, "", -1)
		t.Term = strings.Replace(t.Term, p2, "", -1)
	}
}

func (t *Term) reformat() {
	// Performs more complicated formatting steps
	if strings.Contains(t.Term, "(") == true || strings.Contains(t.Term, ")") == true {
		t.sliceTerm("(", ")")
	}
	if strings.Contains(t.Term, "\"") == true {
		t.sliceTerm("\"", "\"")
	}
	if strings.Contains(t.Term, "/") == true {
		// Subset longer side of slash
		idx := strings.Index(t.Term, "/")
		if idx <= len(t.Term)/2 {
			t.Term = t.Term[idx+1:]
		} else if idx <= len(t.Term)-1 {
			t.Term = t.Term[:idx]
		}
	}
	if strings.Contains(t.Term, "&") == true {
		// Replace ampersand and add spaces if needed
		idx := strings.Index(t.Term, "&")
		if idx > 0 && idx < len(t.Term)-1 {
			if t.Term[idx+1] != ' ' {
				// Check second space first so index remains accurate
				t.Term = t.Term[:idx+1] + " " + t.Term[idx+1:]
			}
			if t.Term[idx-1] != ' ' {
				t.Term = t.Term[:idx] + " " + t.Term[idx:]
			}
			t.Term = strings.Replace(t.Term, "&", "and", 1)
		} else {
			t.Term = strings.Replace(t.Term, "&", "", -1)
		}
	}
	if strings.Contains(t.Term, "#") == true {
		// Drop symbol and any numbers
		idx := strings.Index(t.Term, "#")
		if idx < len(t.Term)/2 {
			ind := strings.Index(t.Term[idx:], " ") + idx
			t.Term = t.Term[ind+1:]
		} else if idx <= len(t.Term)-1 {
			ind := strings.LastIndex(t.Term, " ")
			if ind < idx {
				idx = ind
			}
			// Keep everything up to space/pound
			t.Term = t.Term[:idx]
		}
	}
}

func (t *Term) removeInfant() {
	// Removes words referring to infancy from term
	if strings.Count(t.Term, " ") >= 1 {
		var buffer strings.Builder
		first := true
		s := strings.Split(t.Term, " ")
		words := "Fetus, Juvenile, Infant"
		for _, i := range s {
			if strings.Contains(words, i) == false {
				if first == false {
					buffer.WriteRune(' ')
				}
				buffer.WriteString(i)
				first = false
			}
		}
		t.Term = buffer.String()
	}
}

func (t *Term) checkCertainty() {
	// Sets t.Status if term is unknown or hybrid
	unk := "uncertainEntry"
	hyb := "hybrid"
	l := strings.ToLower(t.Term)
	if strings.Contains(l, "?") == true || strings.Contains(l, "unknown") == true || containsWithSpace(l, "not") == true {
		t.Status = unk
	} else if strings.Contains(l, "hybrid") == true || containsWithSpace(l, "x") == true || containsWithSpace(l, "mix") == true {
		t.Status = hyb
	}
}

func (t *Term) filter() {
	// Filters input query
	query := t.Queries[0]
	if len(query) >= 3 {
		r := regexp.MustCompile(` +`)
		// Replace extra spaces and convert to title case
		t.Term = r.ReplaceAllString(query, " ")
		t.checkCertainty()
		if len(t.Status) == 0 {
			// Convert to title case after checking for ? and x
			t.Term = strarray.TitleCase(t.Term)
			t.removeInfant()
			t.reformat()
			t.checkRunes()
			if len(t.Status) == 0 && len(t.Term) < 3 {
				t.Status = "tooShort"
			}
		}
	} else {
		t.Status = "tooShort"
	}
}
