// Defines extract functions

package main

import (
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"path"
	"regexp"
	"strings"
	"unicode"
)

func (t *term) checkPunctuation() {
	// Returns false if term contains puntuation
	for _, i := range []rune(t.term) {
		if i != '.' && unicode.IsPunct(i) == true {
			t.status = "punctuation"
			break
		}
	}
}

func (t *term) sliceTerm(p1, p2 string) {
	// Removes item from between 2 puntuation marks
	var ind int
	idx := strings.Index(t.term, p1)
	if p1 == p2 {
		ind = strings.LastIndex(t.term, p2)
	} else {
		ind = strings.Index(t.term, p2)
	}
	if idx < ind {
		// Drop item between punctuation
		if ind == len(t.term)-1 {
			t.term = t.term[:idx]
		} else if idx == 0 {
			t.term = t.term[ind+1:]
		} else {
			t.term = t.term[:idx] + t.term[ind+1:]
		}
	} else {
		// Remove puntuation
		t.term = strings.Replace(t.term, p1, "", -1)
		t.term = strings.Replace(t.term, p2, "", -1)
	}
}

func (t *term) compareSlice(s []string, e string) {
	// Sets t.status to e if element in s is in term
	for _, i := range s {
		if strings.Contains(t.term, i) == true {
			t.status = e
			break
		}
	}
}

func (t *term) reformat() {
	// Performs more complicated fortting steps
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

func (t *term) filter() {
	// Filters input query
	query := t.queries[0]
	if len(query) >= 3 {
		r := regexp.MustCompile(` +`)
		// Replace extra spaces and convert to title case
		t.term = r.ReplaceAllString(strings.Title(query), " ")
		t.compareSlice([]string{"?", "not", "unknown"}, "uncertainEntry")
		if len(t.status) == 0 {
			t.compareSlice([]string{" x", "mix ", " mix", "hybrid"}, "hybrid")
			if len(t.status) == 0 {
				// Reformat before filtering for numbers since it might drop number content
				t.reformat()
				t.compareSlice([]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}, "numberContent")
				if len(t.status) == 0 {
					t.checkPunctuation()
				}
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
				t.filter()
				// Append terms with no fail reason to pass; else append to fail
				if len(t.status) == 0 {
					pass = append(pass, []string{t.queries[0], t.term})
				} else {
					fail = append(fail, []string{t.queries[0], t.term, t.status})
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
	pass, fail := filterTerms(*infile, *column)
	fmt.Printf("\tSuccessfully formatted %d entries.\n\t%d entries failed formatting.", len(pass), len(fail))
	iotools.WriteToCSV(*outfile, "Query,SearchTerm\n", pass)
	iotools.WriteToCSV(misses, "Query,SearchTerm,Reason\n", fail)
}
