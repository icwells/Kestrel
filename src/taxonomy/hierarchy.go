// Stores taxonomy hierarchy data for checkTaxonomies

package taxonomy

import (
	"strings"
)

type Hierarchy struct {
	header   map[string]int
	levels   []string
	parents  map[string]string
	children map[string]string
	phylum   map[string]string
	class    map[string]string
	order    map[string]string
	family   map[string]string
	genus    map[string]string
	species  map[string]string
}

func emptyHierarchy() *Hierarchy {
	// Initializes empty taxonomy hierarchy
	h := new(Hierarchy)
	h.levels = []string{"Species", "Genus", "Family", "Order", "Class", "Phylum", "Kingdom"}
	h.parents = map[string]string{"Phylum": "Kingdom",
		"Class":   "Phylum",
		"Order":   "Class",
		"Family":  "Order",
		"Genus":   "Family",
		"Species": "Genus",
	}
	h.children = make(map[string]string)
	for k, v := range h.parents {
		// Reverse parents map
		h.children[v] = k
	}
	h.phylum = make(map[string]string)
	h.class = make(map[string]string)
	h.order = make(map[string]string)
	h.family = make(map[string]string)
	h.genus = make(map[string]string)
	h.species = make(map[string]string)
	return h
}

func NewHierarchy(taxa []*Taxonomy) *Hierarchy {
	// Initializes new taxonomy hierarchy
	h := emptyHierarchy()
	h.setHierarchy(taxa)
	return h
}

func (h *Hierarchy) FillTaxonomy(t *Taxonomy) {
	// Replaces NAs with value from hierarchy
	if strings.ToUpper(t.Genus) == "NA" {
		t.Genus = h.species[t.Species]
	}
	if strings.ToUpper(t.Family) == "NA" {
		t.Family = h.genus[t.Genus]
	}
	if strings.ToUpper(t.Order) == "NA" {
		t.Order = h.family[t.Family]
	}
	if strings.ToUpper(t.Class) == "NA" {
		t.Class = h.order[t.Order]
	}
	if strings.ToUpper(t.Phylum) == "NA" {
		t.Phylum = h.class[t.Class]
	}
	if strings.ToUpper(t.Kingdom) == "NA" {
		t.Kingdom = h.phylum[t.Phylum]
	}
	t.CountNAs()
}

func (h *Hierarchy) setParent(level, parent, child string) {
	// Stores child as key in level map with parent as key (i.e. stores parent level for each child level)
	if child != "NA" && parent != "NA" {
		switch level {
		case "Phylum":
			if _, ex := h.phylum[child]; ex == false {
				h.phylum[child] = parent
			}
		case "Class":
			if _, ex := h.class[child]; ex == false {
				h.class[child] = parent
			}
		case "Order":
			if _, ex := h.order[child]; ex == false {
				h.order[child] = parent
			}
		case "Family":
			if _, ex := h.family[child]; ex == false {
				h.family[child] = parent
			}
		case "Genus":
			if _, ex := h.genus[child]; ex == false {
				h.genus[child] = parent
			}
		case "Species":
			if _, ex := h.species[child]; ex == false {
				h.species[child] = parent
			}
		}
	}
}

func (h *Hierarchy) AddTaxonomy(t *Taxonomy) {
	// Adds individual taxa to hierarchy
	if _, ex := h.phylum[t.Phylum]; !ex {
		h.setParent("Phylum", t.Kingdom, t.Phylum)
	}
	if _, ex := h.class[t.Class]; !ex {
		h.setParent("Class", t.Phylum, t.Class)
	}
	if _, ex := h.order[t.Order]; !ex {
		h.setParent("Order", t.Class, t.Order)
	}
	if _, ex := h.family[t.Family]; !ex {
		h.setParent("Family", t.Order, t.Family)
	}
	if _, ex := h.genus[t.Genus]; !ex {
		h.setParent("Genus", t.Family, t.Genus)
	}
	if _, ex := h.species[t.Species]; !ex {
		h.setParent("Species", t.Genus, t.Species)
	}
}

func (h *Hierarchy) setHierarchy(taxa []*Taxonomy) {
	// Stores corpus in hierarchy
	for _, v := range taxa {
		if v.Nas == 0 {
			h.AddTaxonomy(v)
		}
	}
}
