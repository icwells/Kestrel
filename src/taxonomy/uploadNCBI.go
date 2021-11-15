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
						u.ids[id] = newRank(id, "", name, "")
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

func (u *uploader) printNCBI() {
	// Prints complete ncbi taxonomies to file
	fmt.Println("\tPrinting NCBI records...")
	var res [][]string
	for _, i := range u.taxa {
		res = append(res, i.Slice(i.ID, ""))
	}
	iotools.WriteToCSV("test.csv", "Species...", res)
}

func (u *uploader) loadNCBI() {
	// Uploads NCBI table and formats data into sql database
	u.logger.Println("Reading NCBI taxonomies...")
	u.setNCBIcitations()
	u.setNCBInames()
	for i := range u.YieldNCBI(u.ncbi["nodes"]) {
		id := i[0]
		if level, ex := u.ids[id]; ex {
			level.level = strings.ToLower(i[2])
			level.parent = i[1]
			if level.level == "species" {
				// Store species, genus, and citation
				t := NewTaxonomy()
				t.SetLevel(level.level, level.name)
				t.Genus = level.parent
				t.ID = level.id
				if cit, e := u.citations[id]; e {
					t.Source = cit
				}
				u.taxa = append(u.taxa, t)
				/*} else {
				parents[id] = []string{i[1], i[2]}*/
			}
		}
	}
	u.setLevelIDs()
	u.fillTaxonomies("NCBI")
	u.printNCBI()
	/*fmt.Println("\tUploading NCBI data...")
	u.db.UploadSlice("Taxonomy", u.res)
	u.db.UploadSlice("Common", u.commontable)*/
}
