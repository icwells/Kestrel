// Returns maps of curated taxonomies

package taxonomy

import (
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/kestrel/src/kestrelutils"
	"path"
	"strings"
)

var CORPUS = "corpus.csv.gz"

func setCorpus(infile string) (map[string]*Taxonomy, map[string]string) {
	// Returns taxonomy and common names maps from given file
	taxa := make(map[string]*Taxonomy)
	common := make(map[string]string)
	kestrelutils.CheckFile(infile)
	rows, header := iotools.ReadFile(infile, true)
	for _, i := range rows {
		t := NewTaxonomy()
		c := t.SpeciesCaps(i[header["SearchTerm"]])
		s := t.SpeciesCaps(i[header["Species"]])
		if len(c) > 0 {
			common[c] = s
		}
		t.Kingdom = strings.TrimSpace(i[header["Kingdom"]])
		t.Phylum = strings.TrimSpace(i[header["Phylum"]])
		t.Class = strings.TrimSpace(i[header["Class"]])
		t.Order = strings.TrimSpace(i[header["Order"]])
		t.Family = strings.TrimSpace(i[header["Family"]])
		t.Genus = strings.TrimSpace(i[header["Genus"]])
		t.Species = s
		t.Source = strings.TrimSpace(i[header["Source"]])
		t.Found = true
		t.CountNAs()
		taxa[s] = t
	}
	return taxa, common
}

func GetCorpus() (map[string]*Taxonomy, map[string]string) {
	// Returns taxonomy and common names maps from corpus file
	infile := kestrelutils.GetAbsPath(CORPUS)
	return setCorpus(infile)
}

func corpusHeader() string {
	// Returns Header for corpus
	var ret strings.Builder
	t := NewTaxonomy()
	ret.WriteString("SearchTerm")
	for _, i := range t.levels {
		ret.WriteByte(',')
		ret.WriteString(t.SpeciesCaps(i))
	}
	ret.WriteByte(',')
	ret.WriteString("Source")
	return ret.String()
}

func FormatCorpus(infile string) {
	// Formats new coprus for later searches
	var res [][]string
	outfile := path.Join(kestrelutils.Getutils(), CORPUS)
	taxa, common := setCorpus(infile)
	for _, v := range taxa {
		// Format taxonomy entries
		v.CheckTaxa()
	}
	for i := 0; i <= 1; i++ {
		// Check hierarchy twice to account for corrected NAs
		h := NewHierarchy(taxa)
		for _, v := range taxa {
			h.FillTaxonomy(v)
		}
	}
	for k, v := range common {
		if t, ex := taxa[v]; ex {
			if t.Nas == 0 {
				// Write common names and associated taxonomy and remove from map
				res = append(res, []string{k, t.String()})
			}
			delete(taxa, k)
		}
	}
	for _, v := range taxa {
		if v.Nas == 0 {
			// Write remaining taxa without commoon names
			res = append(res, []string{"", v.String()})
		}
	}
	fmt.Printf("\tWriting new corpus to %s\n", outfile)
	iotools.WriteToCSV(outfile, corpusHeader(), res)
}
