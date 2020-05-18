// Scrapes taxononomies into taxonomy struct

package searchtaxa

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
)

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

type eolstruct struct {
	Entry struct {
		Species string `json:"canonical_form"`
	} `json:"entry"`
	Ancestors []struct {
		ScientificName string `json:"scientificName"`
		TaxonRank      string `json:"taxonRank"`
	} `json:"ancestors"`
}

func (t *taxonomy) scrapeEOL(result []byte, url string) {
	// Scrapes taxonomy from EOL hierarchy entry
	t.source = removeKey(url)
	var j eolstruct
	err := json.Unmarshal(result, &j)
	if err == nil {
		t.setLevel("species", j.Entry.Species)
		for _, a := range j.Ancestors {
			level := t.isLevel(a.TaxonRank)
			if len(level) >= 1 {
				t.setLevel(level, a.ScientificName)
			}
		}
		t.checkTaxa()
	}
}

func (t *taxonomy) scrapeItis(url string) {
	// Scrapes taxonomy info from itis
	t.source = url
	page, err := goquery.NewDocument(url)
	if err == nil {
		found := 0
		page.Find("table").EachWithBreak(func(i int, table *goquery.Selection) bool {
			table.Find("tr").Each(func(i int, tr *goquery.Selection) {
				tr.Find("td").Each(func(j int, td *goquery.Selection) {
					str := td.Text()
					if len(str) > 0 {
						level := t.isLevel(str)
						if len(level) > 0 {
							if level == "species" {
								t.setLevel(level, removeNonBreakingSpaces(td.Next().Text()))
							} else {
								t.setLevel(level, removeNonBreakingSpaces(td.Next().Find("a").Text()))
							}
							found++
						}
					}
				})

			})
			if found == 7 {
				// Break if all levels have been found
				return false
			}
			return true
		})
		t.checkTaxa()
	}
}

type iucnstruct struct {
	// https://mholt.github.io/json-to-go/
	Result []struct {
		Species string `json:"scientific_name"`
		Kingdom string `json:"kingdom"`
		Phylum  string `json:"phylum"`
		Class   string `json:"class"`
		Order   string `json:"order"`
		Family  string `json:"family"`
		Genus   string `json:"genus"`
	} `json:"result"`
}

func (t *taxonomy) scrapeIUCN(result []byte, url string) {
	// Marshalls json array into struct
	t.source = removeKey(url)
	var j iucnstruct
	err := json.Unmarshal(result, &j)
	if err == nil {
		// Map from iucnstruct struct to taxonomy struct
		for _, a := range j.Result {
			//a := j.result[0]
			t.kingdom = a.Kingdom
			t.phylum = a.Phylum
			t.class = a.Class
			t.order = a.Order
			t.family = a.Family
			t.genus = a.Genus
			t.species = a.Species
			t.checkTaxa()
		}
	}
}
