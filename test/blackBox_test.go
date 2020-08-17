// Performs black box tests on kestrel extract and search

package kestrel_test

import (
	"github.com/icwells/go-tools/dataframe"
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
