// Performs black box tests on kestrel extract and search

package kestrel_test

import (
	"github.com/icwells/go-tools/iotools"
	"strings"
	"testing"
)

func readFile(infile string, taxa bool) map[string][]string {
	// Reads in data as map
	ret := make(map[string][]string)
	first := true
	f := iotools.OpenFile(infile)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		if first == false {
			s := strings.Split(strings.TrimSpace(string(scanner.Text())), ",")
			if taxa == true {
				// Use search term as key and ignore sources
				ret[s[1]] = s[1:9]
			} else {
				//Use query as key
				ret[s[0]] = []string{s[1]}
			}
		} else {
			first = false
		}
	}
	return ret
}

func compareFiles(t *testing.T, exp, act string, taxa bool) {
	// Comapres output file to expected
	expected := readFile(exp, taxa)
	actual := readFile(act, taxa)
	if len(expected) != len(actual) {
		t.Errorf("Actual length %d does not equal expected: %d", len(actual), len(expected))
	}
	for k := range actual {
		if _, ex := expected[k]; ex == false {
			t.Errorf("Actual key %s not in expected map.", k)
		} else {
			for idx, i := range actual[k] {
				if i != expected[k][idx] {
					t.Errorf("Actual value for entry %s %s does not equal expected: %s", k, i, expected[k][idx])
				}
			}
		}
	}
}

func TestExtract(t *testing.T) {
	// Tests extract output
	compareFiles(t, "testInput.csv", "extracted.csv", false)
}

func TestSearch(t *testing.T) {
	// Tests search output
	compareFiles(t, "taxonomies.csv", "searchResults.csv", true)
}
