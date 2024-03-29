// Formats and uploads taxonomy databases to MySQL

package taxonomy

import (
	"fmt"
	"github.com/icwells/dbIO"
	"github.com/icwells/kestrel/src/kestrelutils"
	"log"
	"path"
	"strconv"
	"sync"
)

type rank struct {
	id     string
	level  string
	name   string
	parent string
}

func newRank(id, level, name, parent string) *rank {
	// Returns initialized struct
	r := new(rank)
	r.id = id
	r.level = level
	r.name = name
	r.parent = parent
	return r
}

type uploader struct {
	citations   map[string]string
	common      map[string][]string
	commontable [][]string
	count       int
	db          *dbIO.DBIO
	dir         string
	hier        *Hierarchy
	ids         map[string]*rank
	logger      *log.Logger
	names       map[string]string
	ncbi        map[string]string
	proc        int
	res         [][]string
	taxa        []*Taxonomy
	tid         int
}

func newUploader(db *dbIO.DBIO, proc int, logger *log.Logger) *uploader {
	// Returns initialized struct
	u := new(uploader)
	u.citations = make(map[string]string)
	u.common = make(map[string][]string)
	u.db = db
	u.dir = path.Join(kestrelutils.GetLocation(), "databases")
	u.hier = emptyHierarchy()
	u.ids = make(map[string]*rank)
	u.logger = logger
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
	u.ids = make(map[string]*rank)
	u.res = nil
	u.taxa = nil
}

func (u *uploader) storeTaxonomy(wg *sync.WaitGroup, mut *sync.RWMutex, t *Taxonomy, db string) {
	// Fills in missing fields and stores passing taxonomy
	defer wg.Done()
	t.clearids()
	u.hier.FillTaxonomy(t)
	if t.Nas == 0 {
		mut.Lock()
		// Attempt to get existing id
		id, ex := u.names[t.Species]
		if !ex {
			// Upload unique taxonomies
			id = strconv.Itoa(u.tid)
			row := t.Slice(id, db)
			u.res = append(u.res, row)
			u.names[t.Species] = id
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
		mut.Unlock()
	}
}

func (u *uploader) fillTaxonomies(db string) {
	// Fills missing taxonomy fields
	var wg sync.WaitGroup
	var mut sync.RWMutex
	var count int
	u.hier.setHierarchy(u.taxa)
	u.logger.Println("Formatting taxonomies...")
	for _, i := range u.taxa {
		if i.Found {
			wg.Add(1)
			count++
			go u.storeTaxonomy(&wg, &mut, i, db)
			fmt.Printf("\tDispatched %d of %d taxonomies...\r", count, len(u.taxa))
			if count%u.proc == 0 {
				wg.Wait()
			}
		}
	}
	fmt.Println()
	wg.Wait()
}

func UploadDatabases(db *dbIO.DBIO, proc int, logger *log.Logger) {
	// Formats and uploads taxonomy databases to MySQL
	u := newUploader(db, proc, logger)
	u.loadITIS()
}
