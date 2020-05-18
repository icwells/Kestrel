// Defines extract functions

package terms

import (
	"bytes"
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"path"
	"regexp"
	"strings"
	"unicode"
) 

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

func mergeTerms(s []*Term) map[string]*Term {
	// Merges terms which format to same spelling and tries to resolve abbreviations
	ret := make(map[string]*Term)
	for _, i := range s {
		t, ex := ret[i.term]
		if ex {
			i.addQuery(i.query[0])
		} else {
			ret[i.term] = i
			
		}
	}
	fmt.Printf("\tFound %d unique entries from %d total new entries.\n", len(ret), len(s))
	return ret
}

func filterTerms(infile string, c int) ([]*Term, [][]string) {
	// Reads terms from given column and checks formatting
	first := true
	var d string
	var fail [][]string
	var pass []*Term
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
						pass = append(pass, t})
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

func ExtractSearchTerms(infile, outfile string, col int) map[string]*Terms {
	// Extracts and formats input terms
	fmt.Println("\n\tExtracting search terms...")
	kestrelutils.CheckFile(infile)
	dir, _ := path.Split(outfile)
	misses := path.Join(dir, "KestrelRejected.csv")
	pass, fail := filterTerms(infile, col)
	fmt.Printf("\tSuccessfully formatted %d entries.\n\t%d entries failed formatting.\n", len(pass), len(fail))
	iotools.WriteToCSV(misses, "Query,SearchTerm,Reason", fail)
	ret := mergeTerms(pass)
	return ret
}
