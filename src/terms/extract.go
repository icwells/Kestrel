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

type extractor struct {
	col     int
	fail    [][]string
	infile  string
	logger  *log.Logger
	misses  string
	names   []*Term
	outfile string
	speller aspell.Speller
}

func newExtractor(infile, outfile string, col int, logger *log.Logger) *extractor {
	// Returns initialized struct
	kestrelutils.CheckFile(infile)
	e := new(extractor)
	e.col = col
	e.infile = infile
	e.logger = logger
	e.outfile = outfile
	e.speller, _ = aspell.NewSpeller(map[string]string{"lang": "en_US"})
	dir, _ := path.Split(e.outfile)
	e.misses = path.Join(dir, "KestrelRejected.csv")
	return e
}

func (e *extractor) mergeTerms() map[string]*Term {
	// Merges terms which format to same spelling and tries to resolve abbreviations
	ret := make(map[string]*Term)
	for _, i := range e.names {
		if _, ex := ret[i.Term]; ex {
			i.AddQuery(i.Queries[0])
		} else {
			ret[i.Term] = i
		}
	}
	e.logger.Printf("Found %d unique entries from %d total new entries.\n", len(ret), len(e.names))
	return ret
}

func (e *extractor) filterTerms() {
	// Reads terms from given column and checks formatting
	first := true
	var d string
	f := iotools.OpenFile(e.infile)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(string(scanner.Text()))
		if !first {
			var query string
			if e.col >= 0 {
				s := strings.Split(line, d)
				if len(s) > e.col {
					query = s[e.col]
				}
			} else {
				query = line
			}
			if query != "" {
				t := NewTerm(query)
				if len(t.Queries) >= 1 {
					t.filter(e.speller)
					// Append terms with no fail reason to pass; else append to fail
					if len(t.Status) == 0 {
						e.names = append(e.names, t)
					} else {
						e.fail = append(e.fail, []string{t.Queries[0], t.Term, t.Status})
					}
				}
			}
		} else {
			d, _ = iotools.GetDelim(line)
			first = false
		}
	}
}

func ExtractSearchTerms(infile, outfile string, col int, logger *log.Logger) map[string]*Term {
	// Extracts and formats input terms
	e := newExtractor(infile, outfile, col, logger)
	e.filterTerms()
	e.logger.Printf("Successfully formatted %d entries.", len(e.names))
	e.logger.Printf("%d entries failed formatting.", len(e.fail))
	if len(e.fail) > 0 {
		iotools.WriteToCSV(e.misses, "Query,SearchTerm,Reason", e.fail)
	}
	return e.mergeTerms()
}
