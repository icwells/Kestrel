// Formats and uploads taxonomy databases to MySQL

package taxonomy

import (
	"fmt"
	"github.com/icwells/dbIO"
	"github.com/icwells/kestrel/src/kestrelutils"
	"math"
	"path"
	"sync"
	"strconv"
)

type uploader struct {
	common map[string][]string
	count  int
	db     *dbIO.DBIO
	gbif   string
	hier   *Hierarchy
	ids    map[string]string
	itis   string
	names  map[string]string
	ncbi   string
	res    [][]string
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
	u.hier = emptyHierarchy()
	u.ids = make(map[string]string)
	u.itis = "ITIS"
	u.names = make(map[string]string)
	//u.ncbi = path.Join(dir,
	u.tid = 1
	return u
}

func (u *uploader) clear() {
	// Empties taxa slice and common map between datasets
	u.common = make(map[string][]string)
	u.ids = make(map[string]string)
	u.res = nil
	u.taxa = nil
}

func (u *uploader) getDenominator() int {
	// Returns denominator for subsetting upload slice (size in bytes / 16Mb)
	max := 10000000.0
	size := 0
	for _, i := range u.res {
		for _, j := range i {
			size += len([]byte(j))
		}
	}
	return int(math.Ceil(float64(size*8) / max))
}

func (u *uploader) uploadTable(table string) {
	// Uploads patient entries to db
	l := len(u.res)
	if l > 0 {
		den := u.getDenominator()
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
			vals, ln := dbIO.FormatSlice(u.res[ind:end])
			fmt.Println(vals[:500])
			u.db.UpdateDB(table, vals, ln)
			ind = ind + idx
		}
	}
}

func (u *uploader) storeTaxonomy(wg *sync.WaitGroup, mut *sync.RWMutex, t *Taxonomy, db string) {
	// Fills in missing fields and stores passing taxonomy
	defer wg.Done()
	u.hier.FillTaxonomy(t)
	if t.Nas == 0 {
		id := strconv.Itoa(u.tid)
		row := t.Slice(id, db)
		u.res = append(u.res, row)
		u.tid++
		// Store found names
		mut.Lock()
		u.names[t.Species] = id
		mut.Unlock()
	}
}

func (u *uploader) setTaxonomy(wg *sync.WaitGroup, mut *sync.RWMutex, t *Taxonomy) {
	// Replaces rank ids with names
	defer wg.Done()
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
							mut.Lock()
							u.hier.AddTaxonomy(t)
							mut.Unlock()
							u.count++
							t.Found = true
						}
					}
				}
			}
		}
	}
}

func (u *uploader) fillTaxonomies() {
	// Merges taxonomy and ids and fills missing fields
	var wg sync.WaitGroup
	var mut sync.RWMutex
	var count int
	proc := 500
	fmt.Println("\tFilling taxonomies...")
	for _, i := range u.taxa {
		// Fill in taxonomy
		wg.Add(1)
		count++
		go u.setTaxonomy(&wg, &mut, i)
		fmt.Printf("\tDispatched %d of %d taxonomies...\r", count, len(u.taxa))
		if count%proc == 0 {
			wg.Wait()
		}
	}
	fmt.Println()
	wg.Wait()
	count = 0
	fmt.Println("\tFormatting taxonomies...")
	for _, i := range u.taxa {
		if i.Found {
			wg.Add(1)
			count++
			go u.storeTaxonomy(&wg, &mut, i, "GBIF")
			fmt.Printf("\tDispatched %d of %d taxonomies...\r", count, u.count)
			if count%proc == 0 {
				wg.Wait()
			}
		}
	}
	fmt.Println()
	wg.Wait()
}

func UploadDatabases(db *dbIO.DBIO) {
	// Formats and uploads taxonomy databases to MySQL
	u := newUploader(db)
	u.loadGBIF()
	u.clear()
	//u.loadITIS()
	//u.loadNCBI()
}
