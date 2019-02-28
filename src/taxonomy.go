// Defines taxonomy struct and methods

package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

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
	levels  []string
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
	t.levels = []string{"kingdom", "phylum", "class", "order", "family", "genus", "species"}
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

func (t *taxonomy) copyTaxonomy(x taxonomy) {
	// Deep copies x to t
	t.kingdom = x.kingdom
	t.phylum = x.phylum
	t.class = x.class
	t.order = x.order
	t.family = x.family
	t.genus = x.genus
	t.species = x.species
	t.source = x.source
	t.found = x.found
	t.nas = x.nas
}

func (t *taxonomy) countNAs() {
	// Rechecks nas
	nas := 0
	for _, i := range []string{t.kingdom, t.phylum, t.class, t.order, t.family, t.genus, t.species} {
		if strings.ToUpper(i) == "NA" {
			nas++
		}
	}
	t.nas = nas
}

func (t *taxonomy) checkLevel(l string, sp bool) string {
	// Returns formatted name
	if strings.ToUpper(l) != "NA" {
		l = strings.Replace(l, ",", "", -1)
		if sp == false {
			if strings.Contains(l, " ") == true {
				l = strings.Split(l, " ")[0]
			}
			l = titleCase(l)
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
	} else {
		// Standardize NAs
		l = strings.ToUpper(l)
	}
	return l
}

func (t *taxonomy) checkTaxa() {
	// Checks formatting
	t.countNAs()
	if t.nas <= 2 && strings.ToUpper(t.genus) != "NA" {
		t.found = true
		if strings.ToLower(t.kingdom) == "metazoa" {
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
	value = strings.TrimSpace(value)
	if strings.Contains(value, "[") == false && strings.ToUpper(value) != "NA" && len(value) > 1 {
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
		case "scientific_name":
			t.species = value
		}
	}
}

func (t *taxonomy) isLevel(s string) string {
	// Returns formatted string if s is a taxonomic level
	s = strings.TrimSpace(strings.ToLower(strings.Replace(s, ":", "", -1)))
	for _, i := range t.levels {
		if i == s {
			return s
		}
	}
	return ""
}

func (t *taxonomy) scrapeWiki(url string) {
	// Marshalls html taxonomy into struct
	t.source = url
	page, err := goquery.NewDocument(url)
	if err == nil {
		page.Find("td").Each(func(i int, s *goquery.Selection) {
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
			}
		})
		t.checkTaxa()
	}
}

func (t *taxonomy) scrapeNCBI(url string) {
	// Scrapes taxonomy form NCBI efetch results
	t.source = removeKey(url)
	page, err := goquery.NewDocument(url)
	if err == nil {
		taxa := page.Find("Taxon")
		// Get species name
		n := taxa.Find("ScientificName").First()
		r := taxa.Find("Rank").First()
		level := t.isLevel(r.Text())
		if len(level) > 0 {
			t.setLevel(level, n.Text())
		}
		lineage := taxa.Find("LineageEx")
		lineage.Find("Taxon").Each(func(i int, s *goquery.Selection) {
			r = s.Find("Rank")
			level = t.isLevel(r.Text())
			if len(level) > 0 {
				t.setLevel(level, s.Find("ScientificName").Text())
			}
		})
		t.checkTaxa()
	}
}

func (t *taxonomy) scrapeEOL(url string) {
	// Scrapes taxonomy from EOL hierarchy entry
	t.source = removeKey(url)
	page, err := goquery.NewDocument(url)
	if err == nil {
		found := 0
		page.Find("ancestor").EachWithBreak(func(i int, r *goquery.Selection) bool {
			level := t.isLevel(r.Find("taxonRank").Text())
			if len(level) > 0 {
				t.setLevel(level, r.Find("scientificName").Text())
				found++
				if found == 7 {
					// Break if all levels have been found
					return true
				}
			}
			return false
		})
		entry := page.Find("entry")
		// Store canonical species name
		t.setLevel("species", entry.Find("canonical-form").Text())
		t.checkTaxa()
	}
}

func (t *taxonomy) scrapeItis(url string) {
	// Scrapes taxonomy info from itis
	t.source = url
	page, err := goquery.NewDocument(url)
	if err == nil {
		found := 0
		page.Find("tr").EachWithBreak(func(i int, tr *goquery.Selection) bool {
			tr.Find("td").Each(func(j int, td *goquery.Selection) {
				str := td.Text()
				if len(str) > 0 {
					level := t.isLevel(str)
					if len(level) > 0 {
						t.setLevel(level, td.Next().Find("a").Text())
						found++
					}
				}
			})
			if found == 7 {
				// Break if all levels have been found
				return true
			}
			return false
		})
		t.checkTaxa()
		//fmt.Println(t.String())
	}
}

type jsa struct {
	res map[string]string `json:"result"`
}

func (t *taxonomy) scrapeIUCN(result []byte, url string) {
	// Marshalls json array into struct
	t.source = removeKey(url)
	var a jsa
	a.res = make(map[string]string)
	err := json.Unmarshal(result, &a)
	if err == nil {
		// Map from anonymous struct to taxonomy struct
		fmt.Println(string(result))
		for k, v := range a.res {
			fmt.Println(k, v)
			level := t.isLevel(k)
			if len(level) > 1 && len(v) > 2 {
				t.setLevel(level, v)
			}
		}
		//fmt.Println(t.String())
		t.checkTaxa()
	}
}
