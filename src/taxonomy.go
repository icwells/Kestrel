// Defines taxonomy struct and methods

package main

import (
	"encoding/json"
	"io"
	"strings"
)

type apis struct {
	ncbi   string
	wiki   string
	iucn   string
	eol    string
	search string
	pages  string
	hier   string
	format string
}

func newAPIs() apis {
	// Returns api struct
	var a apis
	a.ncbi = "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/"
	a.iucn = "http://apiv3.iucnredlist.org/api/v3/species/"
	a.wiki = "https://en.wikipedia.org/wiki/"
	a.eol = "http://eol.org/api/"
	a.search = "search/1.0."
	a.pages = "pages/1.0."
	a.hier = "hierarchy_entries/1.0."
	a.format = "xml"
	return a
}

type taxonomy struct {
	kingdom string
	phylum  string
	class   string
	order   string
	family  string
	genus   string
	species string
	source  string
	found   bool
	nas     int
}

func newTaxonomy() taxonomy {
	// Initializes taxonomy struct
	var t taxonomy
	t.kingdom = "NA"
	t.phylum = "NA"
	t.class = "NA"
	t.order = "NA"
	t.family = "NA"
	t.genus = "NA"
	t.species = "NA"
	t.source = "NA"
	t.found = false
	t.nas = 7
	return t
}

func (t *taxonomy) String() string {
	// Returns formatted string without source
	var ret []string
	for _, i := range []string{t.kingdom, t.phylum, t.class, t.order, t.family, t.genus, t.species} {
		ret = append(ret, i)
	}
	return strings.Join(ret, ",")
}

func (t *taxonomy) countNAs() {
	// Rechecks nas
	for _, i := range []string{t.kingdom, t.phylum, t.class, t.order, t.family, t.genus, t.species} {
		if i != "NA" {
			t.nas--
		}
	}
}

func (t *taxonomy) checkLevel(l string, sp bool) string {
	// Returns formatted name
	if l != "NA" {
		l = strings.TrimSpace(strings.Replace(l, ",", "", -1))
		if sp == false && strings.Contains(l, " ") == true {
			l = strings.Split(l, " ")[0]
			l = strings.Title(l)
		} else {
			// Get binomial with proper capitalization
			if strings.Contains(l, " ") == false {
				l = t.genus + " " + strings.ToLower(l)
			} else {
				s := strings.Split(l, " ")
				l = strings.Title(s[0]) + " " + strings.ToLower(s[1])
			}
		}
	}
	return l
}

func (t *taxonomy) checkTaxa() {
	// Checks formatting
	t.countNAs()
	if t.nas <= 2 && t.genus != "NA" {
		t.found = true
		if t.kingdom == "Metazoa" {
			// Correct NCBI kingdom
			t.kingdom = "Animalia"
		} else {
			t.kingdom = t.checkLevel(t.kingdom, false)
		}
		t.phylum = t.checkLevel(t.phylum, false)
		t.class = t.checkLevel(t.class, false)
		t.order = t.checkLevel(t.order, false)
		t.family = t.checkLevel(t.family, false)
		t.genus = t.checkLevel(t.genus, false)
		t.species = t.checkLevel(t.species, true)
	}
}

func (t *taxonomy) scrapeWiki(result io.Reader, url string) {
	// Marshalls html taxonomy into struct

}

func (t *taxonomy) scrapeIUCN(result io.Reader, url string) {
	// Marshalls json array into struct
	a := struct {
		result struct {
			kingdom, phylum, class, order, family, genus, scientific_name string
		}
	}{}
	json.NewDecoder(result).Decode(a)
	// Map from anonymous struct to taxonomy struct
	t.kingdom = a.result.kingdom
	t.phylum = a.result.phylum
	t.class = a.result.class
	t.order = a.result.order
	t.family = a.result.family
	t.genus = a.result.genus
	t.species = a.result.scientific_name
	t.checkTaxa()
}
