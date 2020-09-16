// Contains functions for each taxonomy source

package taxonomy

import (
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"strings"
)

func (u *uploader) splitName(n string) (string, string) {
	// Splits citation from name if needed
	c := "NA"
	if strings.Count(n, " ") > 1 {
		s := strings.Split(n, " ")
		n = strings.Join(s[:2], " ")
		c = strings.Join(s[2:], " ")
	}
	return n, strings.Replace(c, ",", "", -1)
}

func (u *uploader) loadGBIF() {
	// Uploads GBIF table and formats data into sql database
	fmt.Println("\tReading GBIF taxonomies...")
	reader, _ := iotools.YieldFile(u.gbif, false)
	for i := range reader {
		if strings.ToUpper(i[4]) == "ACCEPTED" {
			rank := strings.ToLower(i[5])
			if rank == "species" {
				// Store species with ids for ranks
				var sp string
				t := NewTaxonomy()
				sp, t.Source = u.splitName(i[18])
				t.SetLevel("species", sp)
				t.ID = i[0]
				for idx, id := range i[10:16] {
					if id != `\N` {
						t.SetLevel(t.levels[idx], id)
					}
				}
				u.taxa = append(u.taxa, t)
			} else {
				name := i[18]
				if strings.Contains(name, " ") {
					name = strings.Split(name, " ")[0]
				}
				u.ids[i[0]] = name
			}
		}
	}
	u.fillTaxonomies("GBIF")
	fmt.Println("\tUploading GBIF data...")
	u.uploadTable("Taxonomy", u.res)
}

/*func (u *uploader) loadITIS() {
	// Uploads ITIS table and formats data into sql database
	fmt.Println("\tReading ITIS taxonomies...")
}*/

func (u *uploader) YieldNCBI(infile string) <-chan []string {
	ch := make(chan []string)
	d := "|"
	go func() {
		f := iotools.OpenFile(infile)
		defer f.Close()
		input := iotools.GetScanner(f)
		for input.Scan() {
			var s []string
			line := strings.TrimSpace(string(input.Text()))
			for _, i := range strings.Split(line, d) {
				s = append(s, strings.TrimSpace(i))
			}
			ch <- s
		}
		close(ch)
	}()
	return ch
}

func (u *uploader) setNCBIcitations() {
	// Stores ncbi nodes in ids map
	for i := range u.YieldNCBI(u.ncbi["citations"]) {
		if len(i) >= 6 && len(i[6]) > 0 && len(i[1]) > 0 {
			for _, k := range strings.Split(i[6], " ") {
				// Taxa_id: citation
				u.citations[k] = i[1]
			}
		}
	}
}

func (u *uploader) setNCBInames() {
	// Stores ncbi nodes in ids map
	for i := range u.YieldNCBI(u.ncbi["names"]) {
		if len(i) >= 4 {
			id := i[0]
			if i[3] == "scientific name" {
				if name, _ := u.splitName(i[1]); !strings.Contains(name, ".") {
					if strings.Count(name, " ") == 0 || len(strings.Split(name, " ")[1]) > 3 {
						u.ids[id] = name
					}
				}
			} else if i[3] == "common name" {
				if _, ex := u.common[id]; !ex {
					u.common[id] = []string{}
				}
				u.common[id] = append(u.common[id], i[1])
			}
		}
	}
}

func (u *uploader) setNCBIlevels(parents map[string]string) {
	// Stores ids for taxonomic levels
	for _, i := range u.taxa {
		if v, ex := parents[i.Genus]; ex {
			i.Family = v
			if v, ex := parents[i.Family]; ex {
				i.Order = v
				if v, ex := parents[i.Order]; ex {
					i.Class = v
					if v, ex := parents[i.Class]; ex {
						i.Phylum = v
						if v, ex := parents[i.Phylum]; ex {
							i.Kingdom = v
						}
					}
				}
			}
		}
	}
}

func (u *uploader) loadNCBI() {
	// Uploads NCBI table and formats data into sql database
	parents := make(map[string]string)
	fmt.Println("\tReading NCBI taxonomies...")
	u.setNCBIcitations()
	u.setNCBInames()
	for i := range u.YieldNCBI(u.ncbi["nodes"]) {
		id := i[0]
		if i[2] == "species" {
			if name, ex := u.ids[id]; ex {
				// Store species, genus, and citation
				t := NewTaxonomy()
				t.SetLevel("species", name)
				t.Genus = i[1]
				if cit, e := u.citations[id]; e {
					t.Source = cit
				}
				u.taxa = append(u.taxa, t)
			}
		} else {
			parents[id] = i[1]
		}
	}
	u.setNCBIlevels(parents)
	u.fillTaxonomies("NCBI")
	fmt.Println("\tUploading NCBI data...")
	u.uploadTable("Taxonomy", u.res)
	u.uploadTable("Common", u.commontable)
}
