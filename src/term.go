// Defines term struct and methods

package main

import (
	"strings"
)

type term struct {
	queries  []string
	term     string
	status   string
	taxonomy taxonomy
	sources  map[string]string
}

func newTerm(query string) term {
	// Returns initialized term
	var t term
	if len(query) >= 1 {
		t.addQuery(query)
	}
	t.taxonomy = newTaxonomy()
	t.sources = make(map[string]string)
	return t
}

func percentDecode(term string) string {
	// Removes percent encoding from web search
	term = strings.Replace(term, "%20", " ", -1)
	return strings.Replace(term, "%27", "'", -1)
}

func (t *term) String() string {
	// Returns formatted string
	var ret []string
	ret = append(ret, percentDecode(t.term))
	ret = append(ret, t.taxonomy.String())
	// Append url or NA for each source
	for _, i := range []string{"IUCN", "NCBI", "WIKI", "EOL", "ITIS"} {
		s, ex := t.sources[i]
		if ex == true {
			ret = append(ret, s)
		} else {
			ret = append(ret, "NA")
		}
	}
	return strings.Join(ret, ",")
}

func (t *term) addQuery(query string) {
	// Appends to query slice
	t.queries = append(t.queries, query)
}
