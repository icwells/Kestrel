// Performs black box tests on kestrel extract and search

package kestrel_test

import (
	"github.com/icwells/go-tools/dataframe"
	"testing"
)

func TestSearch(t *testing.T) {
	// Tests search output
	exp, _ := dataframe.FromFile("taxonomies.csv", 1)
	act, _ := dataframe.FromFile("searchResults.csv", 1)
	act.DeleteColumn("Source")
	act.DeleteColumn("Confirmed")
	if err := exp.Compare(act); err != nil {
		t.Error(err)
	}
}
