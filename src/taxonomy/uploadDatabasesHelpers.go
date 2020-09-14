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
	u.fillTaxonomies()
	fmt.Println("\tUploading GBIF data...")
	u.uploadTable("Taxonomy")
}

/*func (u *uploader) loadITIS() {
	// Uploads ITIS table and formats data into sql database
	fmt.Println("\tReading ITIS taxonomies...")
}

func (u *uploader) loadNCBI() {
	// Uploads NCBI table and formats data into sql database
	fmt.Println("\tReading NCBI taxonomies...")
}*/
