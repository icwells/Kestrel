// Defines extract functions

package terms

import (
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/kestrel/src/kestrelutils"
	"github.com/trustmaster/go-aspell"
	"log"
	"path"
	"strings"
	"unicode"
)

func containsWithSpace(l, target string) bool {
	// Returns true if target is in l and sperated by spaces/term boundary
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

func mergeTerms(s []*Term, logger *log.Logger) map[string]*Term {
	// Merges terms which format to same spelling and tries to resolve abbreviations
	ret := make(map[string]*Term)
	for _, i := range s {
		if _, ex := ret[i.Term]; ex {
			i.AddQuery(i.Queries[0])
		} else {
			ret[i.Term] = i
		}
	}
	logger.Printf("Found %d unique entries from %d total new entries.\n", len(ret), len(s))
	return ret
}

func filterTerms(infile string, c int) ([]*Term, [][]string) {
	// Reads terms from given column and checks formatting
	first := true
	var d string
	var fail [][]string
	var pass []*Term
	speller, _ := aspell.NewSpeller(map[string]string{"lang": "en_US"})
	f := iotools.OpenFile(infile)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(string(scanner.Text()))
		if !first {
			var query string
			if c >= 0 {
				s := strings.Split(line, d)
				if len(s) > c {
					query = s[c]
				}
			} else {
				query = line
			}
			if query != "" {
				t := NewTerm(query)
				if len(t.Queries) >= 1 {
					t.filter(speller)
					// Append terms with no fail reason to pass; else append to fail
					if len(t.Status) == 0 {
						pass = append(pass, t)
					} else {
						fail = append(fail, []string{t.Queries[0], t.Term, t.Status})
					}
				}
			}
		} else {
			d, _ = iotools.GetDelim(line)
			first = false
		}
	}
	return pass, fail
}

func ExtractSearchTerms(infile, outfile string, col int, logger *log.Logger) map[string]*Term {
	// Extracts and formats input terms
	kestrelutils.CheckFile(infile)
	dir, _ := path.Split(outfile)
	misses := path.Join(dir, "KestrelRejected.csv")
	pass, fail := filterTerms(infile, col)
	logger.Printf("Successfully formatted %d entries.", len(pass))
	logger.Printf("%d entries failed formatting.", len(fail))
	if len(fail) > 0 {
		iotools.WriteToCSV(misses, "Query,SearchTerm,Reason", fail)
	}
	ret := mergeTerms(pass, logger)
	return ret
}
