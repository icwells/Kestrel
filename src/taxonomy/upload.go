// Formats and uploads taxonomy databases to MySQL

package taxonomy

import (
	"fmt"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/kestrel/src/kestrelutils"
	"path"
	"strconv"
	"strings"
	"sync"
)

type uploader struct {
	citations   map[string]string
	common      map[string][]string
	commontable [][]string
	count       int
	db          *dbIO.DBIO
	dir         string
	gbif        string
	hier        *Hierarchy
	ids         map[string]string
	itis        string
	names       map[string]string
	ncbi        map[string]string
	proc        int
	res         [][]string
	taxa        []*Taxonomy
	tid         int
}

func newUploader(db *dbIO.DBIO, proc int) *uploader {
	// Returns initialized struct
	u := new(uploader)
	u.citations = make(map[string]string)
	u.common = make(map[string][]string)
	u.db = db
	u.dir = path.Join(kestrelutils.GetLocation(), "databases")
	u.gbif = path.Join(u.dir, "backbone-current-simple.txt.gz")
	u.hier = emptyHierarchy()
	u.ids = make(map[string]string)
	u.itis = "ITIS"
	u.names = make(map[string]string)
	u.proc = proc
	u.tid = 1
	u.setNCBIfiles()
	return u
}

func (u *uploader) setNCBIfiles() {
	// Stores ncbi names
	u.ncbi = make(map[string]string)
	u.ncbi["citations"] = path.Join(u.dir, "citations.dmp")
	u.ncbi["names"] = path.Join(u.dir, "names.dmp")
	u.ncbi["nodes"] = path.Join(u.dir, "nodes.dmp")
}

func (u *uploader) clear() {
	// Empties taxa slice and common map between datasets
	u.citations = make(map[string]string)
	u.common = make(map[string][]string)
	u.commontable = nil
	u.count = 0
	u.ids = make(map[string]string)
	u.res = nil
	u.taxa = nil
}

func (u *uploader) storeTaxonomy(wg *sync.WaitGroup, mut *sync.RWMutex, t *Taxonomy, db string) {
	// Fills in missing fields and stores passing taxonomy
	defer wg.Done()
	u.hier.FillTaxonomy(t)
	if t.Nas == 0 {
		mut.Lock()
		id, ex := u.names[t.Species]
		if !ex {
			// Upload unique taxonomies
			id = strconv.Itoa(u.tid)
			row := t.Slice(id, db)
			u.res = append(u.res, row)
			u.tid++
		}
		if v, ex := u.common[t.ID]; ex {
			for _, i := range v {
				if _, ex := u.names[i]; !ex {
					// Append common names with new/existing id and store to avoid duplicates
					u.commontable = append(u.commontable, []string{id, i})
					u.names[i] = id
				}
			}
		}
		u.names[t.Species] = id
		mut.Unlock()
	}
}

func (u *uploader) setTaxonomy(wg *sync.WaitGroup, mut *sync.RWMutex, t *Taxonomy) {
	// Replaces rank ids with names
	defer wg.Done()
	if v, ex := u.ids[t.Kingdom]; ex {
		if strings.ToLower(v) == "metazoa" {
			v = "Animalia"
		}
		t.Kingdom = v
		if v, ex = u.ids[t.Phylum]; ex {
			t.Phylum = v
			if v, ex = u.ids[t.Class]; ex {
				t.Class = v
				if v, ex = u.ids[t.Order]; ex {
					t.Order = v
					if v, ex = u.ids[t.Family]; ex {
						t.Family = v
						if v, ex = u.ids[t.Genus]; ex {
							t.Genus = v
							mut.Lock()
							u.hier.AddTaxonomy(t)
							u.count++
							mut.Unlock()
							t.Found = true
						}
					}
				}
			}
		}
	}
}

func (u *uploader) fillTaxonomies(db string) {
	// Merges taxonomy and ids and fills missing fields
	var wg sync.WaitGroup
	var mut sync.RWMutex
	var count int
	fmt.Println("\tFilling taxonomies...")
	for _, i := range u.taxa {
		// Fill in taxonomy
		wg.Add(1)
		count++
		go u.setTaxonomy(&wg, &mut, i)
		fmt.Printf("\tDispatched %d of %d taxonomies...\r", count, len(u.taxa))
		if count%u.proc == 0 {
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
			go u.storeTaxonomy(&wg, &mut, i, db)
			fmt.Printf("\tDispatched %d of %d taxonomies...\r", count, u.count)
			if count%u.proc == 0 {
				wg.Wait()
			}
		}
	}
	fmt.Println()
	wg.Wait()
}

func UploadDatabases(db *dbIO.DBIO, proc int) {
	// Formats and uploads taxonomy databases to MySQL
	u := newUploader(db, proc)
	if iotools.Exists(u.ncbi["nodes"]) {
		u.loadNCBI()
		u.clear()
	}
	if iotools.Exists(u.gbif) {
		u.loadGBIF()
		u.clear()
	}
	//u.loadITIS()
	//os.Remove(u.dir)
}
