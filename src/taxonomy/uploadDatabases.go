// Formats and uploads taxonomy databases to MySQL

package taxonomy

import (
	"fmt"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/kestrel/src/kestrelutils"
	"math"
	"path"
	"strconv"
	"strings"
)

type uploader struct {
	common map[string][]string
	db     *dbIO.DBIO
	gbif   string
	ids    map[string]string
	itis   string
	names  map[string]string
	ncbi   string
	taxa   []*Taxonomy
	tid    int
}

func newUploader(db *dbIO.DBIO) *uploader {
	// Returns initialized struct
	dir := path.Join(kestrelutils.GetLocation(), "databases")
	u := new(uploader)
	u.common = make(map[string][]string)
	u.db = db
	u.gbif = path.Join(dir, "backbone-current-simple.txt.gz")
	u.ids = make(map[string]string)
	u.itis = "ITIS"
	u.names = make(map[string]string)
	//u.ncbi = path.Join(dir,
	u.tid = 1
	return u
}

func (u *uploader) getDenominator(list [][]string) int {
	// Returns denominator for subsetting upload slice (size in bytes / 16Mb)
	max := 10000000.0
	size := 0
	for _, i := range list {
		for _, j := range i {
			size += len([]byte(j))
		}
	}
	return int(math.Ceil(float64(size*8) / max))
}

func (u *uploader) uploadTable(table string, list [][]string) {
	// Uploads patient entries to db
	l := len(list)
	if l > 0 {
		den := u.getDenominator(list)
		if den <= 1 {
			// Upload slice at once
			vals, l := dbIO.FormatSlice(list)
			u.db.UpdateDB(table, vals, l)
		} else {
			// Upload in chunks
			var end int
			idx := l / den
			ind := 0
			for i := 0; i < den; i++ {
				if ind+idx > l {
					// Get last less than idx rows
					end = l
				} else {
					end = ind + idx
				}
				vals, ln := dbIO.FormatSlice(list[ind:end])
				u.db.UpdateDB(table, vals, ln)
				ind = ind + idx
			}
		}
	}
}

func (u *uploader) setTaxonomy(t *Taxonomy) {
	// Replaces rank ids with names
	if v, ex := u.ids[t.Kingdom]; ex {
		t.Kingdom = v
		if v, ex = u.ids[t.Phylum]; ex {
			t.Phylum = v
			if v, ex := u.ids[t.Class]; ex {
				t.Class = v
				if v, ex := u.ids[t.Order]; ex {
					t.Order = v
					if v, ex := u.ids[t.Family]; ex {
						t.Family = v
						if v, ex := u.ids[t.Genus]; ex {
							t.Genus = v
							t.Found = true
						}
					}
				}
			}
		}
	}
}

func (u *uploader) loadGBIF() {
	// Uploads GBIF table and formats data into sql database
	var res [][]string
	fmt.Println("\tReading GBIF taxonomies...")
	reader, _ := iotools.YieldFile(u.gbif, false)
	for i := range reader {
		rank := strings.ToLower(i[5])
		if rank == "species" {
			// Store species with ids for ranks
			t := NewTaxonomy()
			t.SetLevel("species", i[18])
			for idx, id := range i[10:16] {
				t.SetLevel(t.levels[idx], id)
			}
			u.taxa = append(u.taxa, t)
		} else {
			u.ids[i[0]] = i[18]
		}
	}
	for idx, i := range u.taxa {
		// Fill in taxonomy, remove incomplete taxa
		u.setTaxonomy(i)
		if !i.Found {
			if idx == 0 {
				u.taxa = u.taxa[1:]
			} else if idx == len(u.taxa)-1 {
				u.taxa = u.taxa[:idx]
			} else {
				u.taxa = append(u.taxa[:idx], u.taxa[idx+1:]...)
			}
		} else {
			// Record species to avoid multipe entries
			id := strconv.Itoa(u.tid)
			u.names[i.Species] = id
			res = append(res, i.Slice(id, "GBIF"))
		}
	}
	u.uploadTable("Taxonomy", res)
}

/*func (u *uploader) loadITIS() {
	// Uploads ITIS table and formats data into sql database

}

func (u *uploader) loadNCBI() {
	// Uploads NCBI table and formats data into sql database

}*/

func UploadDatabases(db *dbIO.DBIO) {
	// Formats and uploads taxonomy databases to MySQL
	u := newUploader(db)
	u.loadGBIF()
	//u.loadITIS()
	//u.loadNCBI()
}
