// Returns maps of curated taxonomies

package taxonomy

import (
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/kestrel/src/kestrelutils"
	"strings"
)

var CORPUS = "corpus.csv.gz"

func GetCorpus() (map[string]*Taxonomy, map[string]string) {
	// Returns taxonomy and common names maps
	taxa := make(map[string]*Taxonomy)
	common := make(map[string]string)
	infile := kestrelutils.GetAbsPath(CORPUS)
	kestrelutils.CheckFile(infile)
	rows, header := iotools.ReadFile(infile, true)
	for _, i := range rows {
		t := NewTaxonomy()
		c := t.SpeciesCaps(i[header["SearchTerm"]])
		s := t.SpeciesCaps(i[header["Species"]])
		if len(c) > 0 && c != s {
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

func headerString(h map[string]int) string {
	// Returns header as string
	ret := make([]string, len(h))
	for k, v := range h {
		ret[v] = k
	}
	return strings.Join(ret, ",")
}

/*func FormatCorpus(infile string) {
	// Formats corpus for later searching
	var res [][]string
	t := NewTaxonomy()
	outfile := kestrelutils.GetAbsPath(CORPUS)
	kestrelutils.CheckFile(infile)
	rows, header := iotools.ReadFile(infile, true)
	for _, row := range rows {
		var r []string
		c := t.SpeciesCaps(i[header["Query"]])
		s := t.SpeciesCaps(i[header["Species"]])

		for _, i := range row {
			r = append(r, t.SpeciesCaps(i))
		}
		res = append(res, r)
	}
	iotools.WriteToCSV(outfile, headerString(header), res)
}*/
