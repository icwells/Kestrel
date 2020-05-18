// Splits search results file into one with matching search terms and scientific names and one without to streamline manual curration.

package kestrelutils

import (
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"sort"
	"strings"
)

type curated struct {
	taxa map[string][]string
	keys []string
	set  bool
}

func newCurated() curated {
	// Returns initialized curated struct
	var c curated
	c.set = false
	c.taxa = make(map[string][]string)
	return c
}

func (c *curated) setKeys() {
	// Stores keys as slice
	for k := range c.taxa {
		c.keys = append(c.keys, k)
	}
}

func (c *curated) loadTaxa(infile string) {
	// Reads taxa map
	first := true
	f := iotools.OpenFile(infile)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		if first == false {
			line := strings.Split(strings.TrimSpace(string(scanner.Text())), ",")
			c.taxa[line[1]] = line[2:]
		} else {
			first = false
		}
	}
	c.setKeys()
	c.set = true
}

func (c *curated) getTaxonomy(term string) ([]string, bool) {
	// Returns matching taxonomy if term is in map
	var ret []string
	var pass bool
	ret, pass = c.taxa[term]
	if pass == false {
		// Only perform fuzzy search if there is no literal match
		matches := fuzzy.RankFindFold(term, c.keys)
		if len(matches) > 0 {
			sort.Sort(matches)
			if matches[0].Distance >= 0 || matches[0].Distance <= 1 {
				// Accept 0 or 1 transposition
				ret = c.taxa[matches[0].Target]
				pass = true
			}
		}
	}
	return ret, pass
}

func checkTaxonomyResults(infile string, hier hierarchy, taxa curated) (string, [][]string, [][]string) {
	// Identifies records with matching search terms and scientific names
	var header, d string
	var hits, miss [][]string
	var h map[string]int
	first := true
	f := iotools.OpenFile(infile)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(string(scanner.Text()))
		if first == false {
			var pass bool
			s := strings.Split(line, d)
			term := strings.TrimSpace(s[h["SearchTerm"]])
			species := strings.TrimSpace(s[h["Species"]])
			if taxa.set == true {
				row, pass := taxa.getTaxonomy(term)
				if pass == true {
					// Replace with currated taxonomy
					s = append(s[:2], row...)
				}
			}
			if pass == false {
				// Compare search term and species
				score := fuzzy.RankMatchFold(term, species)
				if score >= 0 && score <= 1 {
					pass = true
				}
			}
			// Fill NAs in taxonomy
			s = hier.checkHierarchy(s)
			if pass == true {
				hits = append(hits, s)
			} else {
				miss = append(miss, s)
			}
		} else {
			header = line
			d = iotools.GetDelim(line)
			h = getHeader(strings.Split(line, d))
			first = false
		}
	}
	return header, hits, miss
}

func getOutfiles(name string) (string, string) {
	// Returns formatted output names
	if strings.Contains(name, ".") == true {
		// Remove extension
		name = name[:strings.Index(name, ".")]
	}
	return name + ".passed.csv", name + ".failed.csv"
}

func checkResults() {
	// Checks scientific names in search results
	hier := newHierarchy()
	hier.setLevels(*infile)
	taxa := newCurated()
	if *taxafile != "nil" {
		checkFile(*taxafile)
		// Add to taxonomy hierarchy
		hier.setLevels(*taxafile)
		taxa.loadTaxa(*taxafile)
	}
	pass, fail := getOutfiles(*outfile)
	header, hits, misses := checkTaxonomyResults(*infile, hier, taxa)
	fmt.Println("\tWriting output files...")
	iotools.WriteToCSV(pass, header, hits)
	iotools.WriteToCSV(fail, header, misses)
}
