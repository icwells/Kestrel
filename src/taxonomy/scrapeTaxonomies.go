// Scrapes taxononomies into taxonomy struct

package taxonomy

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/icwells/kestrel/src/kestrelutils"
	"strings"
)

func (t *Taxonomy) ScrapeWiki(url string) {
	// Marshalls html taxonomy into struct
	t.Source = url
	page, err := goquery.NewDocument(url)
	if err == nil {
		page.Find("td").Each(func(i int, s *goquery.Selection) {
			level := t.IsLevel(s.Text(), false)
			if len(level) > 0 {
				n := s.Next()
				name := n.Find("i").Text()
				if len(name) == 0 {
					name = n.Find("a").Text()
				}
				t.SetLevel(level, name)
			}
		})
		t.CheckTaxa()
	}
}

func (t *Taxonomy) ScrapeWikiSpecies(url string) {
	// Marshalls html taxonomy into struct
	t.Source = url
	page, err := goquery.NewDocument(url)
	if err == nil {
		page.Find("p").Each(func(i int, s *goquery.Selection) {
			p := s.Text()
			// check first word of paragraph
			if idx := strings.Index(p, ":"); idx > 0 {
				first := p[:idx]
				if first == "Superregnum" || first == "Familia" {
					for _, i := range strings.Split(p, "\n") {
						row := strings.Split(i, ":")
						if len(row) > 1 {
							if level := t.IsLevel(strings.TrimSpace(row[0]), true); level != "" {
								name := strings.TrimSpace(row[1])
								t.SetLevel(level, name)
							}
						}
					}
				}
			}
		})
		t.CheckTaxa()
	}
}

func (t *Taxonomy) ScrapeAnimalDiversityWeb(url string) {
	// Scrapes html taxonomy into struct
	t.Source = url
	page, err := goquery.NewDocument(url)
	if err == nil {
		page.Find("ul").Each(func(i int, sel *goquery.Selection) {
			if cl, ex := sel.Attr("class"); ex && cl == "unstyled" {
				sel.Find("li").Each(func(j int, s *goquery.Selection) {
					sp := s.Find("span")
					if cl, ex := sp.Attr("class"); ex && cl == "rank" {
						if l := sp.Text(); len(l) > 7 {
							if level := t.ContainsLevel(l[:7]); level != "" {
								name := sp.NextFiltered("a").Text()
								if !strings.Contains(name, ":") {
									t.SetLevel(level, name)
								}
							}
						}
					}
				})
			}
		})
		t.CheckTaxa()
	}
}

func (t *Taxonomy) ScrapeNCBI(url string) {
	// Scrapes taxonomy form NCBI efetch results
	t.Source = kestrelutils.RemoveKey(url)
	page, err := goquery.NewDocument(url)
	if err == nil {
		taxa := page.Find("Taxon")
		// Get species name
		n := taxa.Find("ScientificName").First()
		r := taxa.Find("Rank").First()
		level := t.IsLevel(r.Text(), false)
		if len(level) > 0 {
			t.SetLevel(level, n.Text())
		}
		lineage := taxa.Find("LineageEx")
		lineage.Find("Taxon").Each(func(i int, s *goquery.Selection) {
			r = s.Find("Rank")
			level = t.IsLevel(r.Text(), false)
			if len(level) > 0 {
				t.SetLevel(level, s.Find("ScientificName").Text())
			}
		})
		t.CheckTaxa()
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

func (t *Taxonomy) ScrapeEOL(result []byte, url string) {
	// Scrapes taxonomy from EOL hierarchy entry
	t.Source = kestrelutils.RemoveKey(url)
	var j eolstruct
	err := json.Unmarshal(result, &j)
	if err == nil {
		t.SetLevel("species", j.Entry.Species)
		for _, a := range j.Ancestors {
			level := t.IsLevel(a.TaxonRank, false)
			if len(level) >= 1 {
				t.SetLevel(level, a.ScientificName)
			}
		}
		t.CheckTaxa()
	}
}

func (t *Taxonomy) ScrapeItis(url string) {
	// Scrapes taxonomy info from itis
	t.Source = url
	page, err := goquery.NewDocument(url)
	if err == nil {
		found := 0
		page.Find("table").EachWithBreak(func(i int, table *goquery.Selection) bool {
			table.Find("tr").Each(func(i int, tr *goquery.Selection) {
				tr.Find("td").Each(func(j int, td *goquery.Selection) {
					str := td.Text()
					if len(str) > 0 {
						level := t.IsLevel(str, false)
						if len(level) > 0 {
							if level == "species" {
								t.SetLevel(level, kestrelutils.RemoveNonBreakingSpaces(td.Next().Text()))
							} else {
								t.SetLevel(level, kestrelutils.RemoveNonBreakingSpaces(td.Next().Find("a").Text()))
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
		t.CheckTaxa()
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

func (t *Taxonomy) ScrapeIUCN(result []byte, url string) {
	// Marshalls json array into struct
	t.Source = kestrelutils.RemoveKey(url)
	var j iucnstruct
	err := json.Unmarshal(result, &j)
	if err == nil {
		// Map from iucnstruct struct to taxonomy struct
		for _, a := range j.Result {
			t.Kingdom = a.Kingdom
			t.Phylum = a.Phylum
			t.Class = a.Class
			t.Order = a.Order
			t.Family = a.Family
			t.Genus = a.Genus
			t.Species = a.Species
			t.CheckTaxa()
		}
	}
}
