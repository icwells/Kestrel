// Sorts and re-uploads itis taxa

package taxonomy

import (
	"fmt"
	"os"
	"strings"
)

func (u *uploader) itisKingdoms() {
	// Returns itis kingdom map
	for _, i := range u.db.GetTable("kingdoms") {
		u.ids[i[0]] = newRank(i[0], "kingdom", i[1], "")
	}
}

func (u *uploader) itisRanks() map[string]map[string]string {
	// Returns itis ranks stored by kingdom id
	ranks := make(map[string]map[string]string)
	for _, i := range u.db.GetColumns("taxon_unit_types", []string{"kingdom_id", "rank_id", "rank_name"}) {
		// Store ranks by rank id and rank ids by kingdom id
		if _, ex := ranks[i[0]]; !ex {
			ranks[i[0]] = make(map[string]string)
		}
		ranks[i[0]][i[1]] = strings.ToLower(i[2])
	}
	return ranks
}

func (u *uploader) setLevelIDs() {
	// Stores ids for taxonomic levels
	u.logger.Println("Sorting IDs...")
	for _, i := range u.taxa {
		id := i.Genus
		v, ex := u.ids[i.Genus]
		for ex {
			id = v.id
			if ex {
				i.SetLevel(v.level, v.name)
				if v.level == "kingdom" {
					i.Found = true
					break
				}
			}
			/*if v.level == "genus" {
				fmt.Println(i.Genus)
			}*/
			v, ex = u.ids[v.parent]
			if v.id == id {
				break
			}
		}
	}
}

func (u *uploader) setids() {
	// Loads itis ids
	species := "species"
	ranks := u.itisRanks()
	for _, i := range u.db.GetRows("taxonomic_units", "name_usage", "valid", "tsn,parent_tsn,kingdom_id,rank_id,complete_name") {
		id := i[0]
		kid := i[2]
		rid := i[3]
		name := i[4]
		if k, ex := ranks[kid]; ex {
			if level, e := k[rid]; e {
				// Only store if rank can be identified
				if level == species {
					if _, ex := u.names[name]; !ex {
						// Store species, and genus
						t := NewTaxonomy()
						t.SetLevel(species, name)
						t.Genus = i[1]
						//t.Kingdom = kingdoms[kid]
						t.ID = id
						if cit, e := u.citations[id]; e {
							t.Source = cit
						}
						u.taxa = append(u.taxa, t)
					}
				} else {
					u.ids[id] = newRank(id, level, name, i[1])
				}
			}
		}
	}
}

func (u *uploader) setITIScitations() {
	// Stores strippedauthor table in citations map
	for _, i := range u.db.GetTable("strippedauthor") {
		u.citations[i[0]] = i[1]
	}
}

func (u *uploader) setcommon(i []string) {
	// Stores common names
	tsn := i[0]
	if _, ex := u.names[tsn]; !ex {
		if _, ex := u.common[tsn]; !ex {
			u.common[tsn] = []string{}
		}
	}
	u.common[tsn] = append(u.common[tsn], i[1])
}

func (u *uploader) getcommon() {
	// Stores common names by tsn
	table := "vernaculars"
	column := "language"
	columns := "tsn,vernacular_name"
	for _, language := range []string{"English", "unspecified"} {
		for _, i := range u.db.GetRows(table, column, language, columns) {
			u.setcommon(i)
		}
	}
}

func (u *uploader) loadITIS() {
	// Uploads ITIS table and formats data into sql database
	fmt.Println()
	u.logger.Println("Reading ITIS taxonomies...")
	// Close upload connection
	//itis, err := dbIO.Connect(u.db.Host, "ITIS", u.db.User, u.db.Password)
	_, err := u.db.DB.Exec("USE ITIS;")
	if err != nil {
		u.logger.Printf("[Error] Cannot connect to ITIS database: %v\n", err)
		os.Exit(100)
	}
	u.getcommon()
	u.setITIScitations()
	u.itisKingdoms()
	u.setids()
	u.setLevelIDs()
	u.fillTaxonomies("ITIS")
	// Revert to taxonomy database
	u.db.DB.Exec(fmt.Sprintf("USE %s;", u.db.Database))
	u.logger.Println("Uploading ITIS data...")
	u.db.UploadSlice("Taxonomy", u.res)
	u.db.UploadSlice("Common", u.commontable)
}
