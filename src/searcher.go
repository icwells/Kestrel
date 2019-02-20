// Defines searcher struct and methods

package main

import (
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/go-tools/strarray"
	"path"
	"strings"
)

type searcher struct {
	outfile string
	missed  string
	keys    map[string]string
	done    strarray.Set
	terms   map[string]*term
	misses  []string
	urls    apis
	matches int
}

func (s *searcher) assignKey(line string) {
	// Assigns individual api key to struct
	l := strings.Split(line, "=")
	if len(l) == 2 {
		s.keys[strings.TrimSpace(l[0])] = strings.TrimSpace(l[1])
	}
}

func (s *searcher) apiKeys() {
	// Reads api keys from file
	infile := "API.txt"
	checkFile(infile)
	fmt.Println("\tReading API keys from file...")
	f := iotools.OpenFile(infile)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		line := string(scanner.Text())
		if line[0] != '#' {
			s.assignKey(line)
		}
	}
}

func (s *searcher) checkOutput(outfile, header string) {
	// Reads in completed searches
	if iotools.Exists(outfile) == true {
		var d string
		first := true
		fmt.Printf("\tReading previous output from %s\n", outfile)
		out := iotools.OpenFile(outfile)
		defer out.Close()
		scanner := iotools.GetScanner(out)
		for scanner.Scan() {
			line := string(scanner.Text())
			if first == false {
				l := strings.Split(line, d)
				// Store queries (distinct lines)
				s.done.Add(strings.TrimSpace(l[0]))
			} else {
				d = iotools.GetDelim(line)
				first = false
			}
		}
		fmt.Printf("\tFound %d completed entries.\n", s.done.Length())
	} else {
		fmt.Println("\tGenerating new output file...")
		out := iotools.CreateFile(outfile)
		defer out.Close()
		out.WriteString(header + "\n")
	}
}

func newSearcher() searcher {
	// Reads api keys and existing output and initializes maps
	var s searcher
	s.outfile = *outfile
	dir, _ := path.Split(s.outfile)
	s.missed = path.Join(dir, "KestrelMissed.csv")
	s.keys = make(map[string]string)
	s.done = strarray.NewSet()
	s.terms = make(map[string]*term)
	s.urls = newAPIs()
	s.apiKeys()
	s.checkOutput(s.outfile, "Query,SearchTerm,Kingdom,Phylum,Class,Order,Family,Genus,Species,IUCN,NCBI,Wikipedia,EOL,ITIS")
	s.checkOutput(s.missed, "Query,SearchTerm,Reason")
	return s
}

func percentEncode(term string) string {
	// Percent encodes apostrophes and spaces
	term = strings.Replace(term, " ", "%20", -1)
	return strings.Replace(term, "'", "%27", -1)
}

func (s *searcher) termMap(infile string) {
	// Reads formatted species names
	var d string
	var unique, total int
	first := true
	checkFile(infile)
	fmt.Println("\tReading search terms from file...")
	f := iotools.OpenFile(infile)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(string(scanner.Text()))
		if first == false {
			l := strings.Split(line, d)
			if len(l) >= 2 {
				query := strings.TrimSpace(l[0])
				searchterm := percentEncode(strings.TrimSpace(l[1]))
				if s.done.InSet(query) == false {
					total++
					if _, ex := s.terms[searchterm]; ex == false {
						// Initialize new struct
						unique++
						t := newTerm(query)
						t.term = searchterm
						s.terms[searchterm] = &t
					} else {
						// Add to exisiting struct
						s.terms[searchterm].addQuery(query)
					}
				}
			}
		} else {
			d = iotools.GetDelim(line)
			first = false
		}
	}
	fmt.Printf("\tFound %d unique entries from %d total new entries.\n", unique, total)
}

func (s *searcher) writeMatches(k string) {
	// Appends matches to file
	out := iotools.AppendFile(s.outfile)
	match := s.terms[k].String()
	for _, i := range s.terms[k].queries {
		out.WriteString(fmt.Sprintf("%s,%s\n", i, match))
		s.matches++
	}
}
