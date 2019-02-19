// Defines taxonomy struct and methods

package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
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
	levels	string
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
	t.levels = "kingdom,phylum,class,order,family,genus,species"
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
		if sp == false {
			if strings.Contains(l, " ") == true {
				l = strings.Split(l, " ")[0]
			}
			l = strings.Title(l)
		} else {
			// Get binomial with proper capitalization
			if strings.Contains(l, ".") == true {
				// Remove genus abbreviations
				l = strings.TrimSpace(l[strings.Index(l, ".")+1:])
			}
			if strings.Contains(l, " ") == false {
				l = fmt.Sprintf("%s %s", t.genus, strings.ToLower(l))
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

func (t *taxonomy) setLevel(key, value string) {
	// Sets level denoted by key with value
	switch key {
		case "kingdom":
			t.kingdom = value
		case "phylum":
			t.phylum = value
		case "class":
			t.class = value
		case "order":
			t.order = value
		case "family":
			t.family = value
		case "genus":
			t.genus = value
		case "species":
			t.species = value
	}
	t.nas--
}

func (t *taxonomy) isLevel(s string) string {
	// Returns formatted string if s is a taxonomic level
	var ret string
	s = strings.TrimSpace(strings.ToLower(strings.Replace(s, ":", "", -1)))
	if strings.Contains(t.levels, s) == true {
		ret = s
	}
	return ret
}

func (t *taxonomy) scrapeWiki(url string) {
	// Marshalls html taxonomy into struct
	t.source = url
	page, err := goquery.NewDocument(t.source)
	if err == nil {
		page.Find("td").Each(func (i int, s *goquery.Selection) {
			level := t.isLevel(s.Text())
			if len(level) > 0 {
				var a *goquery.Selection
				n := s.Next()
				if level != "species" {
					a = n.Find("a")
				} else {
					a = n.Find("i")
				}
				t.setLevel(level, a.Text())
				//fmt.Printf("Content of cell %d: %s, %s\n", i, level, a.Text())
			}
		})
		t.checkTaxa()
	}
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
	t.countNAs()
	t.checkTaxa()
}
