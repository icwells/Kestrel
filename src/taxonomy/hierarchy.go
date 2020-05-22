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

func NewHierarchy(taxa map[string]*Taxonomy) *Hierarchy {
	// Initializes new taxonomy hierarchy
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

func (h *Hierarchy) AddTaxonomy(t *Taxonomy) {
	// Adds individual taxa to hierarchy
	if _, ex := h.phylum[t.Phylum]; !ex {
		h.phylum[t.Phylum] = t.Kingdom
	}
	if _, ex := h.class[t.Class]; !ex {
		h.class[t.Class] = t.Phylum
	}
	if _, ex := h.order[t.Order]; !ex {
		h.order[t.Order] = t.Class
	}
	if _, ex := h.family[t.Family]; !ex {
		h.family[t.Family] = t.Order
	}
	if _, ex := h.genus[t.Genus]; !ex {
		h.genus[t.Genus] = t.Family
	}
	if _, ex := h.species[t.Species]; !ex {
		h.species[t.Species] = t.Genus
	}
}

func (h *Hierarchy) setHierarchy(taxa map[string]*Taxonomy) {
	// Stores corpus in hierarchy
	for _, v := range taxa {
		h.AddTaxonomy(v)
	}
}

/*func (h *Hierarchy) getParent(level, name string) string {
	// Returns parent level name for given level
	var parent string
	var ex bool
	name = strarray.TitleCase(name)
	switch level {
	case "Phylum":
		parent, ex = h.phylum[name]
	case "Class":
		parent, ex = h.class[name]
	case "Order":
		parent, ex = h.order[name]
	case "Family":
		parent, ex = h.family[name]
	case "Genus":
		parent, ex = h.genus[name]
	case "Species":
		parent, ex = h.species[name]
	}
	if ex == true {
		return parent
	}
	return ""
}

func (h *Hierarchy) checkHierarchy(s []string) []string {
	// Checks row for NAs and replaces if parent is found in struct
	// Iterate backwards starting from genus to fill multiple empty cells
	for _, level := range h.levels[1:] {
		idx := h.header[level]
		if strings.ToLower(s[idx]) == "na" {
			// Get name and index of child level
			child := h.children[level]
			ind := h.header[child]
			parent := h.getParent(child, s[ind])
			if len(parent) >= 1 {
				s[idx] = parent
			}
		}
	}
	return s
}

func (h *Hierarchy) setParent(level, parent, child string) {
	// Stores child as key in level map with parent as key (i.e. stores parent level for each child level)
	if child != "NA" && parent != "NA" {
		parent = strarray.TitleCase(parent)
		child = strarray.TitleCase(child)
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
}*/
