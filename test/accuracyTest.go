// Performs accuracy test for Kestrel package

package main

import (
	"fmt"
	"github.com/icwells/go-tools/dataframe"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/go-tools/strarray"
	"github.com/icwells/kestrel/src/searchtaxa"
	"github.com/icwells/kestrel/src/terms"
	"github.com/icwells/simpleset"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var (
	infile   = "../utils/corpus.csv.gz"
	outfile  = "searchResults.csv"
	col      = 0
	nocorpus = true
	proc     = 50
)

func formatPercent(a, b int) string {
	// Returns a/b as percent
	return strconv.FormatFloat(float64(a)/float64(b)*100.0, 'f', 3, 64) + "%"
}

func speciesCaps(name string) string {
	// Properly capitalizes species name
	name = strings.TrimSpace(strings.ToLower(name))
	s := strings.Split(name, " ")
	if len(s) > 1 {
		// Save with genus capitalized and species in lower case
		var builder strings.Builder
		builder.WriteString(strarray.TitleCase(s[0]))
		for _, i := range s[1:] {
			builder.WriteByte(' ')
			builder.WriteString(i)
		}
		return builder.String()
	} else {
		return strarray.TitleCase(name)
	}
}

type comparison struct {
	confirmed int
	correct   int
	matches   int
	nas       int
	total     int
}

func (c *comparison) String() string {
	var ret strings.Builder
	ret.WriteString(fmt.Sprintf("\n\tFound matches for %d of %d terms (%s).\n", c.matches, c.total, formatPercent(c.matches, c.total)))
	ret.WriteString(fmt.Sprintf("\tFound %d correct results (%s).\n", c.correct, formatPercent(c.correct, c.matches)))
	ret.WriteString(fmt.Sprintf("\tFound %d NAs (%s).\n", c.nas, formatPercent(c.nas, c.matches)))
	ret.WriteString(fmt.Sprintf("\tFound %d confirmed matches (%s).\n\n", c.confirmed, formatPercent(c.confirmed, c.matches)))
	return ret.String()
}

func compareResults(act, exp *dataframe.Dataframe) {
	// Counts total number of correct, missed, etc.
	c := new(comparison)
	c.total = proc
	c.matches = act.Length()
	for k := range act.Index {
		pass := true
		if conf, _ := act.GetCell(k, "Confirmed"); conf == "yes" {
			c.confirmed++
		} else if sp, _ := act.GetCell(k, "Species"); sp == "NA" {
			c.nas++
		}
		for col := range exp.Header {
			a, _ := act.GetCell(k, col)
			e, _ := exp.GetCell(k, col)
			if a != e {
				fmt.Println(k, col, a, e)
				pass = false
				break
			}
		}
		if pass {
			c.correct++
		}
	}
	fmt.Print(c)
}

func subsetTerms(searchterms map[string]*terms.Term) map[string]*terms.Term {
	// Randomly reduces map to 500 entries
	var keys []string
	for k := range searchterms {
		keys = append(keys, k)
	}
	for len(searchterms) > proc {
		idx := rand.Intn(len(keys))
		delete(searchterms, keys[idx])
		keys = strarray.DeleteSliceIndex(keys, idx)
	}
	return searchterms
}

func setExpected() *dataframe.Dataframe {
	// Reads expected dataframe without duplicate indeces
	var exp [][]string
	set := simpleset.NewStringSet()
	rows, header := iotools.ReadFile(infile, true)
	head := make([]string, len(header))
	for k, v := range header {
		head[v] = k
	}
	for _, i := range rows {
		if len(rows[0]) == 0 {
			// Ignore rows without common names
			break
		}
		if ex, _ := set.InSet(i[col]); !ex {
			exp = append(exp, i)
		}
		set.Add(i[col])
	}
	exp = append([][]string{head}, exp...)
	ret, err := dataframe.FromSlice(exp, col)
	if err != nil {
		panic(err)
	}
	ret.DeleteColumn("Source")
	return ret
}

func main() {
	start := time.Now()
	fmt.Println("\n\tExtracting search terms...")
	searchterms := terms.ExtractSearchTerms(infile, outfile, col)
	fmt.Printf("\tCurrent run time: %v\n", time.Since(start))
	fmt.Println("\n\tSearching for taxonomy matches...")
	searchtaxa.SearchTaxonomies(outfile, subsetTerms(searchterms), proc, nocorpus)
	fmt.Printf("\tFinished. Run time: %v\n\n", time.Since(start))
	fmt.Println("\tComparing output...")
	exp := setExpected()
	act, _ := dataframe.FromFile(outfile, 1)
	act.DeleteColumn("Query")
	act.DeleteColumn("Source")
	compareResults(act, exp)
}
