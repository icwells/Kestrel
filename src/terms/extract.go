// Defines extract functions

package terms

import (
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/kestrel/src/kestrelutils"
	"github.com/trustmaster/go-aspell"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
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
	dir     string
	fail    [][]string
	infile  string
	logger  *log.Logger
	merged  map[string]*Term
	min     float64
	misses  string
	names   []*Term
	outfile string
	script  string
	speller aspell.Speller
}

func newExtractor(infile, outfile string, col int, logger *log.Logger) *extractor {
	// Returns initialized struct
	kestrelutils.CheckFile(infile)
	e := new(extractor)
	e.col = col
	e.dir = path.Join(iotools.GetGOPATH(), "src/github.com/icwells/kestrel/nlp/")
	e.infile = infile
	e.logger = logger
	e.merged = make(map[string]*Term)
	e.min = 0.98
	e.outfile = outfile
	e.script = "namePredictor.py"
	e.speller, _ = aspell.NewSpeller(map[string]string{"lang": "en_US"})
	dir, _ := path.Split(e.outfile)
	e.misses = path.Join(dir, "KestrelRejected.csv")
	return e
}

func (e *extractor) writeTerms(outfile string) {
	// Writes input file for name classifier
	out := iotools.CreateFile(outfile)
	defer out.Close()
	for k := range e.merged {
		out.WriteString(k + "\n")
	}
}

func (e *extractor) getClassifications(infile string) {
	// Reads name classifications from file
	reader, _ := iotools.YieldFile(infile, false)
	for i := range reader {
		if v, ex := e.merged[i[0]]; ex {
			if val, err := strconv.ParseFloat(i[1], 64); err == nil {
				if val >= e.min {
					v.Scientific = true
					if s := strings.Split(i[0], " "); len(s) > 2 {
						v.Term = strings.Join(s[:2], " ")
					}
				}
			}
		}
	}
}

func (e *extractor) classifyTerms() {
	// Calls name classifier and assigns values
	orig, _ := os.Getwd()
	infile, outfile := "names.txt", "results.csv"
	e.logger.Println("Calling scientific name classifier...")
	os.Chdir(e.dir)
	defer os.Chdir(orig)
	e.writeTerms(infile)
	defer os.Remove(infile)
	cmd := exec.Command("python", e.script, infile, outfile)
	if err := cmd.Run(); err != nil {
		e.logger.Printf("Name classifier failed. %v\n", err)
	} else {
		defer os.Remove(outfile)
		e.getClassifications(outfile)
	}
	for _, v := range e.merged {
		if !v.Scientific {
			// Check spelling for common names
			v.checkSpelling(e.speller)
		}
	}
}

func (e *extractor) mergeTerms() {
	// Merges terms which format to same spelling and tries to resolve abbreviations
	for _, i := range e.names {
		if _, ex := e.merged[i.Term]; ex {
			i.AddQuery(i.Queries[0])
		} else {
			e.merged[i.Term] = i
		}
	}
	e.logger.Printf("Found %d unique entries from %d total new entries.\n", len(e.merged), len(e.names))
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
					t.filter()
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
	e.mergeTerms()
	e.classifyTerms()
	return e.merged
}
