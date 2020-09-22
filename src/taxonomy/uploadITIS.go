// Sorts and re-uploads itis taxa

package taxonomy

import (
	"fmt"
	"github.com/icwells/dbIO"
)

func (u *uploader) setids(itis *dbIO.DBIO) {
	// Loads itis ids
	ranks := make(make(map[string]map[string]string)
	kingdoms := make(map[string]string)
	for _, i := range itis.GetTable("kingdoms") {
		kingdoms[i[0]] = i[1]
	}
	for _, i := range itis.GetColumns("taxon_unit_types", []string{"kingdom_id", "rank_id", "rank_name"}) {
		// Store ranks by rank id and rank ids by kingdom id
		if _, ex := ranks[i[0]]; !ex {
			ranks[i[0]] = make(map[string]string)
		}
		ranks[i[0]][i[1]] = i[2]
	}
	for _, i := range itis.GetRows("taxonomic_units", "name_usage", "valid", "tsn,parent_tsn,kingdom_id,rank_id,complete_name") {
		
	}
}

func (u *uploader) setcommon(i []string) {
	// Stores common names
	if _, ex := u.common[i[0]]; !ex {
		u.common[i[0]] = []string{}
	}
	u.common[i[0]] = append(u.common[i[0]], i[1])
}

func (u *uploader) getcommon(itis *dbIO.DBIO) {
	// Stores common names by tsn
	table := "vernaculars"
	column := "language"
	columns := "tsn,vernaculars"
	for _, i := range itis.GetRows(table, column, "English", columns) {
		u.setcommon(i)
	}
	for _, i := range itis.GetRows(table, column, "unspecified", columns) {
		u.setcommon(i)
	}
}

func (u *uploader) loadITIS() {
	// Uploads ITIS table and formats data into sql database
	fmt.Println("\n\tReading ITIS taxonomies...")
	itis := dbIO.NewDBIO(u.db.Host, "ITIS", u.DB.User, u.DB.Password)
	u.getcommon(itis)
}
