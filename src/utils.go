// Heler functions and structs

package main

import (
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"os"
	"strings"
)

type apis struct {
	itis   string
	ncbi   string
	wiki   string
	iucn   string
	eol    string
	search string
	pages  string
	hier   string
}

func newAPIs() apis {
	// Returns api struct
	var a apis
	a.itis = "https://www.itis.gov/"
	a.ncbi = "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/"
	a.iucn = "http://apiv3.iucnredlist.org/api/v3/species/"
	a.wiki = "https://en.wikipedia.org/wiki/"
	a.eol = "http://eol.org/api/"
	a.search = "search/1.0."
	a.pages = "pages/1.0."
	a.hier = "hierarchy_entries/1.0."
	return a
}

//----------------------------------------------------------------------------

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
	if len(query) > 0 {
		t.addQuery(query)
	}
	t.taxonomy = newTaxonomy()
	t.sources = make(map[string]string)
	return t
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

//----------------------------------------------------------------------------

func titleCase(t string) string {
	// Manually converts term to title case (strings.Title is buggy)
	var query []string
	s := strings.Split(t, " ")
	for _, i := range s {
		if len(i) > 1 {
			// Skip stray characters
			query = append(query, strings.ToUpper(string(i[0]))+strings.ToLower(i[1:]))
		}
	}
	return strings.Join(query, " ")
}

func checkFile(infile string) {
	// Makes sure imut file exists
	if iotools.Exists(infile) == false {
		fmt.Printf("\n\t[Error] Input file %s not found. Exiting.\n\n", infile)
		os.Exit(1)
	}
}

func percentEncode(term string) string {
	// Percent encodes apostrophes and spaces
	term = strings.Replace(term, " ", "%20", -1)
	return strings.Replace(term, "'", "%27", -1)
}

func percentDecode(term string) string {
	// Removes percent encoding from web search
	term = strings.Replace(term, "%20", " ", -1)
	return strings.Replace(term, "%27", "'", -1)
}

func removeKey(url string) string {
	// Returns urls with api key removed
	var idx int
	if strings.Contains(url, "&") == true {
		idx = strings.LastIndex(url, "&")
	} else if strings.Contains(url, "?") == true {
		idx = strings.LastIndex(url, "&")
	}
	if idx > 0 {
		url = url[:idx]
	}
	return url
}
