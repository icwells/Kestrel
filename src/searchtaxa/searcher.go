// Defines searcher struct and methods

package searchtaxa

import (
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/kestrel/src/kestrelutils"
	"github.com/icwells/kestrel/src/taxonomy"
	"github.com/icwells/kestrel/src/terms"
	"github.com/icwells/simpleset"
	"path"
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

func newAPIs() *apis {
	// Returns api struct
	a := new(apis)
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

type searcher struct {
	common  map[string]string
	corpus  bool
	done    *simpleset.Set
	fails   int
	hier    *taxonomy.Hierarchy
	keys    map[string]string
	matches int
	missed  string
	outfile string
	service service
	taxa    map[string]*taxonomy.Taxonomy
	terms   map[string]*terms.Term
	urls    *apis
}

func newSearcher(outfile string, searchterms map[string]*terms.Term, nocorpus, test bool) searcher {
	// Reads api keys and existing output and initializes maps
	var s searcher
	s.corpus = !nocorpus
	s.outfile = outfile
	dir, _ := path.Split(s.outfile)
	s.missed = path.Join(dir, "KestrelMissed.csv")
	s.keys = make(map[string]string)
	s.done = simpleset.NewStringSet()
	s.taxa, s.common = taxonomy.GetCorpus()
	s.hier = taxonomy.NewHierarchy(s.taxa)
	s.terms = searchterms
	s.urls = newAPIs()
	if test == false {
		s.service = newService()
		s.apiKeys()
		s.checkOutput(s.outfile, "Query,SearchTerm,Kingdom,Phylum,Class,Order,Family,Genus,Species,Source,Confirmed")
		s.checkOutput(s.missed, "Query,SearchTerm")
	}
	return s
}

func (s *searcher) assignKey(line string) {
	// Assigns individual api key to struct
	l := strings.Split(line, "=")
	if len(l) == 2 {
		s.keys[strings.TrimSpace(l[0])] = strings.TrimSpace(l[1])
	}
}

func (s *searcher) apiKeys() {
	// Reads api keys from file
	infile := kestrelutils.GetAbsPath("API.txt")
	kestrelutils.CheckFile(infile)
	fmt.Println("\tReading API keys from file...")
	f := iotools.OpenFile(infile)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		line := string(scanner.Text())
		if line[0] != '#' {
			s.assignKey(line)
		}
	}
}

func (s *searcher) checkOutput(outfile, header string) {
	// Reads in completed searches
	l := s.done.Length()
	if iotools.Exists(outfile) == true {
		var d string
		first := true
		fmt.Printf("\tReading previous output from %s\n", outfile)
		out := iotools.OpenFile(outfile)
		defer out.Close()
		scanner := iotools.GetScanner(out)
		for scanner.Scan() {
			line := string(scanner.Text())
			if first == false {
				l := strings.Split(line, d)
				// Store queries (distinct lines)
				s.done.Add(strings.TrimSpace(l[0]))
			} else {
				d, _ = iotools.GetDelim(line)
				first = false
			}
		}
		fmt.Printf("\tFound %d completed entries.\n", s.done.Length()-l)
	} else {
		fmt.Println("\tGenerating new output file...")
		out := iotools.CreateFile(outfile)
		defer out.Close()
		out.WriteString(header + "\n")
	}
}

func (s *searcher) writeMisses(k string) {
	// Writes terms with no match to missed file
	out := iotools.AppendFile(s.missed)
	defer out.Close()
	t := kestrelutils.PercentDecode(k)
	for _, i := range s.terms[k].Queries {
		out.WriteString(fmt.Sprintf("%s,%s\n", i, t))
		s.fails++
	}
}

func (s *searcher) writeMatches(k string) {
	// Appends matches to file
	out := iotools.AppendFile(s.outfile)
	defer out.Close()
	match := s.terms[k].String()
	for _, i := range s.terms[k].Queries {
		out.WriteString(fmt.Sprintf("%s,%s\n", i, match))
		s.matches++
	}
}
