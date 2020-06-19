// Performs black box tests on kestrel extract and search

package kestrel_test

import (
	"fmt"
	"github.com/icwells/go-tools/dataframe"
	"strconv"
	"strings"
	"testing"
)

var (
	INFILE = "searchResults.csv"
	EXP    = "../utils/corpus.csv.gz"
)

func expectedTaxa() [][]string {
	return [][]string{
		{"Query", "SearchTerm", "Kingdom", "Phylum", "Class", "Order", "Family", "Genus", "Species"},
		{"Coyote", "Coyote", "Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Canis", "Canis latrans"},
		{"Canis Latrans", "Canis latrans", "Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Canis", "Canis latrans"},
		{"canis lupus", "Canis lupus", "Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Canis", "Canis lupus"},
		{"wolf", "Wolf", "Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Canis", "Canis lupus"},
		{"GRAY WOLF", "Gray wolf", "Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Canis", "Canis lupus"},
		{"GRAY FOX (frank)", "Gray fox", "Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Urocyon", "Urocyon cinereoargenteus"},
		{"Urocyon cinereoargenteus", "Urocyon cinereoargenteus", "Animalia", "Chordata", "Mammalia", "Carnivora", "Canidae", "Urocyon", "Urocyon cinereoargenteus"},
	}
}

func TestSearch(t *testing.T) {
	// Tests search output
	exp, _ := dataframe.FromSlice(expectedTaxa(), 1)
	act, _ := dataframe.FromFile(INFILE, 1)
	act.DeleteColumn("Source")
	act.DeleteColumn("Confirmed")
	if err := exp.Compare(act); err != nil {
		t.Error(err)
	}
}

func formatPercent(a, b int) string {
	// Returns a/b as percent
	return strconv.FormatFloat(float64(a)/float64(b)*100.0, 'f', 3, 64) + "%"
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
	ret.WriteString(fmt.Sprintf("\n\tFound matches for %d of %d terms (%s).", c.matches, c.total, formatPercent(c.matches, c.total)))
	ret.WriteString(fmt.Sprintf("\tFound %d correct results (%s).", c.correct, formatPercent(c.correct, c.matches)))
	ret.WriteString(fmt.Sprintf("\tFound %d NAs (%s).", c.nas, formatPercent(c.nas, c.matches)))
	ret.WriteString(fmt.Sprintf("\tFound %d confirmed matches (%s).", c.confirmed, formatPercent(c.confirmed, c.matches)))
	return ret.String()
}

func (c *comparison) compareResults(t *testing.T, act, exp *dataframe.Dataframe) {
	// Counts total number of correct, missed, etc.
	c.total = exp.Length()
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
				pass = false
				break
			}
		}
		if pass {
			c.correct++
		}
	}

}

func TestFullSearch(t *testing.T) {
	c := new(comparison)
	exp, _ := dataframe.FromFile(EXP, 0)
	exp.DeleteColumn("Source")
	exp.RenameColumn("Common", "SearchTerm")
	act, _ := dataframe.FromFile(INFILE, 1)
	act.DeleteColumn("Source")
	act.DeleteColumn("Confirmed")
	c.compareResults(t, act, exp)
	t.Error(c)
}
