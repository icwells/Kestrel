// Splits search results file into one with matching search terms and scientific names and one without to streamline manual curration.

package main

import (
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"strings"
)

func checkTaxonomyResults(infile string) (string, [][]string, [][]string) {
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
			s := strings.Split(line, d)
			if strings.TrimSpace(s[h["SearchTerm"]]) == strings.TrimSpace(s[h["Species"]]) {
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
	checkFile(*infile)
	pass, fail := getOutfiles(*outfile)
	header, hits, misses := checkTaxonomyResults(*infile)
	fmt.Println("\tWriting output files...")
	iotools.WriteToCSV(pass, header, hits)
	iotools.WriteToCSV(fail, header, misses)
}
