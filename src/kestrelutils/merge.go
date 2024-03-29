// Merges search results with source data

package kestrelutils

import (
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"log"
	"strings"
)

func getHeader(row []string) map[string]int {
	// Returns map of header indeces
	ret := make(map[string]int)
	for idx, i := range row {
		ret[i] = idx
	}
	return ret
}

type taxamerger struct {
	taxa map[string][]string
	nas  []string
}

func newTaxa(infile string, logger *log.Logger) taxamerger {
	// Reads in results as a map of string slices
	var t taxamerger
	t.taxa = make(map[string][]string)
	t.nas = []string{"NA", "NA", "NA", "NA", "NA", "NA", "NA"}
	var d string
	var h map[string]int
	first := true
	logger.Println("Reading search result file...")
	f := iotools.OpenFile(infile)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(string(scanner.Text()))
		if first == false {
			s := strings.Split(line, d)
			// Query name: [taxonomy] (drops search term and urls)
			if len(s) >= h["Species"] {
				t.taxa[s[h["Query"]]] = s[h["Kingdom"] : h["Species"]+1]
				// Additionally store scientific name as key
				t.taxa[s[h["Species"]]] = s[h["Kingdom"] : h["Species"]+1]
			}
		} else {
			d, _ = iotools.GetDelim(line)
			h = getHeader(strings.Split(line, d))
			first = false
		}
	}
	return t
}

func (t *taxamerger) getTaxa(n string) []string {
	// Returns taxonomy for given name
	ret, ex := t.taxa[n]
	if ex != true {
		ret = t.nas
	}
	return ret
}

func (t *taxamerger) mergeTaxonomy(infile string, c int, prepend bool, logger *log.Logger) (string, [][]string) {
	// Returns header and merged results
	first := true
	var ret [][]string
	var d, header string
	fmt.Println("Merging input file with taxonomies...")
	f := iotools.OpenFile(infile)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(string(scanner.Text()))
		if first == false {
			s := strings.Split(line, d)
			if len(s) >= c {
				var row []string
				taxa := t.getTaxa(s[c])
				if prepend == false {
					row = append(s, taxa...)
				} else {
					row = append(taxa, s...)
				}
				ret = append(ret, row)
			}
		} else {
			d, _ = iotools.GetDelim(line)
			s := strings.Split(line, d)
			if prepend == false {
				header = strings.Join(s, ",") + ",Kingdom,Phylum,Class,Order,Family,Genus,ScientificName"
			} else {
				header = "Kingdom,Phylum,Class,Order,Family,Genus,ScientificName," + strings.Join(s, ",")
			}
			first = false
		}
	}
	return header, ret
}

func MergeResults(infile, resfile, outfile string, col int, prepend bool, logger *log.Logger) {
	// Merges search results with source file
	CheckFile(infile)
	CheckFile(resfile)
	taxa := newTaxa(resfile, logger)
	header, results := taxa.mergeTaxonomy(infile, col, prepend, logger)
	logger.Println("Writing output...")
	iotools.WriteToCSV(outfile, header, results)
}
