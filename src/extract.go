// Defines extract functions

package main

import (
	"bytes"
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"path"
	"regexp"
	"strings"
	"unicode"
)

func (t *term) checkRunes() {
	// Removes puntuation and numbers from term
	var name bytes.Buffer
	for _, i := range []rune(t.term) {
		if unicode.IsLetter(i) == true || unicode.IsSpace(i) == true {
			// Remove punctuation and numbers
			name.WriteRune(i)
		} else if i == '.' || i == '-' || i == '\'' {
			name.WriteRune(i)
		}
	}
	t.term = name.String()
	// Double check starting and ending runes for escaped punctuation
	if len(t.term) > 0 {
		if t.term[0] == '.' || t.term[0] == '-' {
			t.term = t.term[1:]
		}
		if t.term[len(t.term)-1] == '-' {
			t.term = t.term[:len(t.term)-1]
		}
	}
}

func (t *term) sliceTerm(p1, p2 string) {
	// Removes item from between 2 puntuation marks
	idx := strings.Index(t.term, p1)
	ind := strings.LastIndex(t.term, p2)
	if idx >= 0 && idx < ind {
		// Drop item between punctuation
		if ind == len(t.term)-1 {
			t.term = t.term[:idx]
		} else if idx == 0 {
			t.term = t.term[ind+1:]
		} else {
			t.term = t.term[:idx] + t.term[ind+1:]
		}
		t.term = strings.TrimSpace(t.term)
	} else {
		// Remove puntuation
		t.term = strings.Replace(t.term, p1, "", -1)
		t.term = strings.Replace(t.term, p2, "", -1)
	}
}

func (t *term) reformat() {
	// Performs more complicated formatting steps
	if strings.Contains(t.term, "(") == true || strings.Contains(t.term, ")") == true {
		t.sliceTerm("(", ")")
	}
	if strings.Contains(t.term, "\"") == true {
		t.sliceTerm("\"", "\"")
	}
	if strings.Contains(t.term, "/") == true {
		// Subset longer side of slash
		idx := strings.Index(t.term, "/")
		if idx <= len(t.term)/2 {
			t.term = t.term[idx+1:]
		} else if idx <= len(t.term)-1 {
			t.term = t.term[:idx]
		}
	}
	if strings.Contains(t.term, "&") == true {
		// Replace ampersand and add spaces if needed
		idx := strings.Index(t.term, "&")
		if idx > 0 && idx < len(t.term)-1 {
			if t.term[idx+1] != ' ' {
				// Check second space first so index remains accurate
				t.term = t.term[:idx+1] + " " + t.term[idx+1:]
			}
			if t.term[idx-1] != ' ' {
				t.term = t.term[:idx] + " " + t.term[idx:]
			}
			t.term = strings.Replace(t.term, "&", "and", 1)
		} else {
			t.term = strings.Replace(t.term, "&", "", -1)
		}
	}
	if strings.Contains(t.term, "#") == true {
		// Drop symbol and any numbers
		idx := strings.Index(t.term, "#")
		if idx < len(t.term)/2 {
			ind := strings.Index(t.term[idx:], " ") + idx
			t.term = t.term[ind+1:]
		} else if idx <= len(t.term)-1 {
			ind := strings.LastIndex(t.term, " ")
			if ind < idx {
				idx = ind
			}
			// Keep everything up to space/pound
			t.term = t.term[:idx]
		}
	}
}

func containsWithSpace(l, target string) bool {
	// Returns true is target is in l and sperated by spaces/term boundary
	var ret bool
	idx := strings.Index(l, target)
	if idx >= 0 {
		next := idx + len(target)
		if idx == 0 {
			if next < len(l) && unicode.IsSpace(rune(l[next])) == true {
				// First word
				ret = true
			}
		} else if next < len(l) {
			if unicode.IsSpace(rune(l[next])) == true && unicode.IsSpace(rune(l[idx-1])) == true {
				ret = true
			}
		} else if unicode.IsSpace(rune(l[idx-1])) == true {
			// Last word
			ret = true
		}
	}
	return ret
}

func (t *term) removeInfant() {
	// Removes words referring to infancy from term
	if strings.Count(t.term, " ") >= 1 {
		var buffer bytes.Buffer
		first := true
		s := strings.Split(t.term, " ")
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
		t.term = buffer.String()
	}
}

func (t *term) checkCertainty() {
	// Sets t.status if term is unknown or hybrid
	unk := "uncertainEntry"
	hyb := "hybrid"
	l := strings.ToLower(t.term)
	if strings.Contains(l, "?") == true || strings.Contains(l, "unknown") == true || containsWithSpace(l, "not") == true {
		t.status = unk
	} else if strings.Contains(l, "hybrid") == true || containsWithSpace(l, "x") == true || containsWithSpace(l, "mix") == true {
		t.status = hyb
	}
}

func (t *term) filter() {
	// Filters input query
	query := t.queries[0]
	if len(query) >= 3 {
		r := regexp.MustCompile(` +`)
		// Replace extra spaces and convert to title case
		t.term = r.ReplaceAllString(query, " ")
		t.checkCertainty()
		if len(t.status) == 0 {
			// Convert to title case after checking for ? and x
			t.term = titleCase(t.term)
			t.removeInfant()
			t.reformat()
			t.checkRunes()
			if len(t.status) == 0 && len(t.term) < 3 {
				t.status = "tooShort"
			}
		}
	} else {
		t.status = "tooShort"
	}
}

func filterTerms(infile string, c int) ([][]string, [][]string) {
	// Reads terms from given column and checks formatting
	first := true
	var d string
	var pass, fail [][]string
	f := iotools.OpenFile(infile)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		line := string(scanner.Text())
		if first == false {
			s := strings.Split(line, d)
			if len(s) > c {
				t := newTerm(s[c])
				if len(t.queries) >= 1 {
					t.filter()
					// Append terms with no fail reason to pass; else append to fail
					if len(t.status) == 0 {
						pass = append(pass, []string{t.queries[0], t.term})
					} else {
						fail = append(fail, []string{t.queries[0], t.term, t.status})
					}
				}
			}
		} else {
			d = iotools.GetDelim(line)
			first = false
		}
	}
	return pass, fail
}

func extractSearchTerms() {
	// Extracts and formats input terms
	checkFile(*infile)
	dir, _ := path.Split(*outfile)
	misses := path.Join(dir, "KestrelRejected.csv")
	pass, fail := filterTerms(*infile, *ecol)
	fmt.Printf("\tSuccessfully formatted %d entries.\n\t%d entries failed formatting.\n", len(pass), len(fail))
	iotools.WriteToCSV(*outfile, "Query,SearchTerm", pass)
	iotools.WriteToCSV(misses, "Query,SearchTerm,Reason", fail)
}
