// Defines term struct

package terms

import (
	"github.com/icwells/go-tools/strarray"
	"github.com/icwells/kestrel/src/kestrelutils"
	"github.com/icwells/kestrel/src/taxonomy"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/trustmaster/go-aspell"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

var MAXDIST = 2

type Term struct {
	Confirmed  bool
	Corrected  string
	Queries    []string
	Scientific bool
	Status     string
	Taxonomy   *taxonomy.Taxonomy
	Term       string
}

func NewTerm(query string) *Term {
	// Returns initialized term
	t := new(Term)
	if len(query) > 0 {
		t.AddQuery(query)
	}
	t.Taxonomy = taxonomy.NewTaxonomy()
	return t
}

func (t *Term) String() string {
	// Returns formatted string
	var ret []string
	ret = append(ret, kestrelutils.PercentDecode(t.Term))
	ret = append(ret, t.Taxonomy.String())
	if t.Confirmed {
		ret = append(ret, "yes")
	} else {
		ret = append(ret, "no")
	}
	return strings.Join(ret, ",")
}

func (t *Term) AddQuery(query string) {
	// Appends to query slice
	t.Queries = append(t.Queries, query)
}

func (t *Term) Confirm() {
	// Sets confirmed to true
	t.Confirmed = true
}

func (t *Term) checkSpelling(speller aspell.Speller) {
	// Stores potential corrected spelling in t.Corrected if word is incorrectly spelled
	var builder strings.Builder
	var pass bool
	for idx, i := range strings.Split(t.Term, " ") {
		match := i
		if !speller.Check(match) {
			matches := fuzzy.RankFindFold(i, speller.Suggest(i))
			if matches.Len() > 0 {
				sort.Sort(matches)
				if matches[0].Distance <= MAXDIST {
					match = matches[0].Target
					pass = true
				}
			}
		}
		if idx > 0 {
			builder.WriteByte(' ')
			builder.WriteString(strings.ToLower(match))
		} else {
			builder.WriteString(strarray.TitleCase(match))
		}
	}
	if pass {
		t.Corrected = t.Taxonomy.SpeciesCaps(builder.String())
	}
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
	} else if strings.Contains(l, "hybrid") == true && !containsWithSpace(l, "hybrida") {
		// Exclude hyrids but keep scientific name 'hybrida'
		t.Status = hyb
	} else if containsWithSpace(l, "x") == true || containsWithSpace(l, "mix") == true {
		t.Status = hyb
	}
}

func (t *Term) speciesCaps() {
	// Properly capitalizes species name
	t.Term = t.Taxonomy.SpeciesCaps(t.Term)
}

func (t *Term) filter() {
	// Filters input query
	short := "tooShort"
	query := t.Queries[0]
	if len(query) >= 3 {
		r := regexp.MustCompile(` +`)
		// Replace extra spaces and convert to title case
		t.Term = r.ReplaceAllString(query, " ")
		t.checkCertainty()
		if len(t.Status) == 0 {
			// Convert to title case after checking for ? and x
			t.speciesCaps()
			t.removeInfant()
			t.reformat()
			t.checkRunes()
			if len(t.Status) == 0 && len(t.Term) < 3 {
				t.Status = short
			}
		}
	} else {
		t.Status = short
	}
}
