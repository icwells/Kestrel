// Formats and uploads taxonomy databases to MySQL

package taxonomy

import (
	"fmt"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/kestrel/src/kestrelutils"
	"github.com/icwells/simpleset"
	"math"
	"path"
)

func sizeOf(list [][]string) int {
	// Returns size of array in bytes
	ret := 0
	for _, i := range list {
		for _, j := range i {
			ret += len([]byte(j))
		}
	}
	return ret * 8
}

func getDenominator(size int) int {
	// Returns denominator for subsetting upload slice (size in bytes / 16Mb)
	max := 10000000.0
	return int(math.Ceil(float64(size) / max))
}

type uploader struct {
	common	[][]string
	db		*dbIO.DBIO
	gbif	string
	ids		map[string]string
	itis	string
	names	*simpleset.Set
	ncbi	string
	sources [][]string
	taxa	*Hierarchy
}

func newUploader(db *dbIO.DBIO) *uploader {
	// Returns initialized struct
	dir := path.Join(kestrelutils.GetLocation(), "databases")
	u := make(uploader)
	u.db = db
	u.gbif = path.Join(dir, "backbone-current-simple.txt.gz")
	u.ids = make(map[string]string)
	u.itis = "ITIS"
	u.names = simpleset.NewStringSet()
	//u.ncbi = path.Join(dir, 
	u.taxa = emptyHierarchy()
	return u
}

func (u *uploader) loadGBIF() {
	// Uploads GBIF table and formats data into sql database
	reader, _ := iotools.YieldFile(u.gbif, false)
	for i := range reader {
		
	}
}

func (u *uploader) loadITIS() {
	// Uploads ITIS table and formats data into sql database

}

func (u *uploader) loadNCBI() {
	// Uploads NCBI table and formats data into sql database

}

func UploadDatabases(db *dbIO.DBIO) {
	// Formats and uploads taxonomy databases to MySQL
	u := newUploader(db)
	u.loadGBIF()
	u.loadITIS()
	u.loadNCBI()
}
