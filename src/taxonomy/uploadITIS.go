// Sorts and re-uploads itis taxa

package taxonomy

import (
	"fmt"
	"os"
	"strings"
)

func (u *uploader) itisKingdoms() map[string]string {
	// Returns itis kingdom map
	ret := make(map[string]string)
	for _, i := range u.db.GetTable("kingdoms") {
		ret[i[0]] = i[1]
	}
	return ret
}

func (u *uploader) itisRanks() map[string]map[string]string {
	// Returns itis ranks stored by kingdom id
	ranks := make(map[string]map[string]string)
	for _, i := range u.db.GetColumns("taxon_unit_types", []string{"kingdom_id", "rank_id", "rank_name"}) {
		// Store ranks by rank id and rank ids by kingdom id
		if _, ex := ranks[i[0]]; !ex {
			ranks[i[0]] = make(map[string]string)
		}
		ranks[i[0]][i[1]] = i[2]
	}
	return ranks
}

func (u *uploader) setLevelIDs(parents map[string]string) {
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
						if i.Kingdom == "NA" {
							if v, ex := parents[i.Phylum]; ex {
								i.Kingdom = v
							}
						}
					}
				}
			}
		}
	}
}

func (u *uploader) setids() {
	// Loads itis ids
	parents := make(map[string]string)
	kingdoms := u.itisKingdoms()
	ranks := u.itisRanks()
	for _, i := range u.db.GetRows("taxonomic_units", "name_usage", "valid", "tsn,parent_tsn,kingdom_id,rank_id,complete_name") {
		id := i[0]
		kid := i[2]
		rid := i[3]
		if k, ex := ranks[kid]; ex {
			if rank, e := k[rid]; e {
				// Only store if rank can be identified
				rank = strings.ToLower(rank)
				if rank == "species" {
					if _, ex := u.names[i[4]]; !ex {
						// Store species, and genus
						t := NewTaxonomy()
						t.SetLevel("species", i[4])
						t.Genus = i[1]
						t.Kingdom = kingdoms[kid]
						t.ID = id
						if cit, e := u.citations[id]; e {
							t.Source = cit
						}
						u.taxa = append(u.taxa, t)
					}
				} else {
					parents[id] = i[1]
					u.ids[id] = i[4]
				}
			}
		}
	}
	u.setLevelIDs(parents)
}

func (u *uploader) setITIScitations() {
	// Stores strippedauthor table in citations map
	for _, i := range u.db.GetTable("strippedauthor") {
		u.citations[i[0]] = i[1]
	}
}

func (u *uploader) setcommon(i []string) {
	// Stores common names
	if _, ex := u.names[i[0]]; !ex {
		if _, ex := u.common[i[0]]; !ex {
			u.common[i[0]] = []string{}
		}
		u.common[i[0]] = append(u.common[i[0]], i[1])
	}
}

func (u *uploader) getcommon() {
	// Stores common names by tsn
	table := "vernaculars"
	column := "language"
	columns := "tsn,vernacular_name"
	for _, i := range u.db.GetRows(table, column, "English", columns) {
		u.setcommon(i)
	}
	for _, i := range u.db.GetRows(table, column, "unspecified", columns) {
		u.setcommon(i)
	}
}

func (u *uploader) loadITIS() {
	// Uploads ITIS table and formats data into sql database
	fmt.Println("\n\tReading ITIS taxonomies...")
	// Close upload connection
	//itis, err := dbIO.Connect(u.db.Host, "ITIS", u.db.User, u.db.Password)
	_, err := u.db.DB.Exec("USE ITIS;")
	if err != nil {
		fmt.Printf("\n\t[Error] Cannot connect to ITIS database: %v\n", err)
		os.Exit(100)
	}
	u.getcommon()
	u.setids()
	u.fillTaxonomies("ITIS")
	// Revert to taxonomy database
	u.db.DB.Exec(fmt.Sprintf("USE %s;", u.db.Database))
	fmt.Println("\tUploading NCBI data...")
	u.db.UploadSlice("Taxonomy", u.res)
	u.db.UploadSlice("Common", u.commontable)
}
