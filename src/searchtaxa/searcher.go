// Defines searcher struct and methods

package searchtaxa

import (
	"fmt"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/kestrel/src/kestrelutils"
	"github.com/icwells/kestrel/src/taxonomy"
	"github.com/icwells/kestrel/src/terms"
	"github.com/icwells/simpleset"
	"log"
	"path"
	"strings"
)

type apis struct {
	adw    string
	eol    string
	hier   string
	itis   string
	iucn   string
	ncbi   string
	pages  string
	search string
	wiki   string
	wksp   string
}

func newAPIs() *apis {
	// Returns api struct
	a := new(apis)
	a.adw = "https://animaldiversity.org/"
	a.eol = "http://eol.org/api/"
	a.hier = "hierarchy_entries/1.0."
	a.itis = "https://www.itis.gov/"
	a.iucn = "http://apiv3.iucnredlist.org/api/v3/species/"
	a.ncbi = "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/"
	a.pages = "pages/1.0."
	a.search = "search/1.0."
	a.wiki = "https://en.wikipedia.org/wiki/"
	a.wksp = "https://species.wikimedia.org/wiki/"
	return a
}

//----------------------------------------------------------------------------

type searcher struct {
	common  map[string]string
	corpus  bool
	db      *dbIO.DBIO
	done    *simpleset.Set
	fails   int
	hier    *taxonomy.Hierarchy
	keys    map[string]string
	logger  *log.Logger
	matches int
	missed  string
	names   []string
	outfile string
	service *service
	taxa    map[string]*taxonomy.Taxonomy
	terms   map[string]*terms.Term
	urls    *apis
}

func newSearcher(db *dbIO.DBIO, logger *log.Logger, outfile string, searchterms map[string]*terms.Term, nocorpus, test bool) searcher {
	// Reads api keys and existing output and initializes maps
	var s searcher
	s.corpus = !nocorpus
	s.db = db
	s.outfile = outfile
	dir, _ := path.Split(s.outfile)
	s.missed = path.Join(dir, "KestrelMissed.csv")
	s.keys = make(map[string]string)
	s.done = simpleset.NewStringSet()
	s.hier = taxonomy.NewHierarchy(s.taxa)
	s.logger = logger
	s.terms = searchterms
	s.urls = newAPIs()
	s.getCorpus()
	if test == false {
		s.service = newService()
		s.apiKeys()
		s.checkOutput(s.outfile, "Query,SearchTerm,Kingdom,Phylum,Class,Order,Family,Genus,Species,Source,Confirmed")
		s.checkOutput(s.missed, "Query,SearchTerm")
	}
	return s
}

func (s *searcher) getCorpus() {
	// Stores common name and taxonomy corpus
	common := make(map[string][]string)
	s.common = make(map[string]string)
	set := simpleset.NewStringSet()
	s.taxa = make(map[string]*taxonomy.Taxonomy)
	for _, i := range s.db.GetTable("Common") {
		if _, ex := common[i[0]]; !ex {
			common[i[0]] = []string{}
		}
		common[i[0]] = append(common[i[0]], i[1])
	}
	for _, i := range s.db.GetTable("Taxonomy") {
		id := i[0]
		t := taxonomy.NewTaxonomy()
		t.Kingdom = i[1]
		t.Phylum = i[2]
		t.Class = i[3]
		t.Order = i[4]
		t.Family = i[5]
		t.Genus = i[6]
		t.Species = i[7]
		t.Source = i[9]
		if i[8] != "" {
			// Add citation if available
			t.Source += ": " + i[8]
		}
		s.taxa[t.Species] = t
		set.Add(t.Species)
		if v, ex := common[id]; ex {
			for _, name := range v {
				s.common[name] = t.Species
				set.Add(name)
			}
		}
	}
	s.names = set.ToStringSlice()
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
	s.logger.Println("Reading API keys from file...")
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
		s.logger.Printf("Reading previous output from %s\n", outfile)
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
		s.logger.Printf("Found %d completed entries.\n", s.done.Length()-l)
	} else {
		s.logger.Println("Generating new output file...")
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
