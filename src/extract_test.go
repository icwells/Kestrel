// Tests funcitons in the extract script

package main

import (
	"testing"
)

type extractinput struct {
	query  string
	term   string
	status string
}

func extractentry(s []string) extractinput {
	// Returns struct from slice input
	var ret extractinput
	ret.query = s[0]
	ret.term = s[1]
	ret.status = s[2]
	return ret
}

func newExtractInput() []extractinput {
	// Returns initialized slice of test data
	var ret []extractinput
	ret = append(ret, extractentry([]string{"FISH (CARDINAL TETRA OR PENCIL FISH)", "Fish", ""}))
	ret = append(ret, extractentry([]string{"PIPING` GUAN ", "Piping Guan", ""}))
	ret = append(ret, extractentry([]string{`SEBA'S  STRIPED FINGERFISH "Sheila"`, "Seba's Striped Fingerfish", ""}))
	ret = append(ret, extractentry([]string{"axolotl-5", "", "numberContent"}))
	ret = append(ret, extractentry([]string{"unknown fish", "", "uncertainEntry"}))
	ret = append(ret, extractentry([]string{"ferret?", "", "uncertainEntry"}))
	ret = append(ret, extractentry([]string{"canine mix", "", "hybrid"}))
	ret = append(ret, extractentry([]string{"corgi x", "", "hybrid"}))
	ret = append(ret, extractentry([]string{"xy", "", "tooShort"}))
	return ret
}

func TestFilter(t *testing.T) {
	exp := newExtractInput()
	for _, e := range exp {
		a := newTerm(e.query)
		a.filter()
		if len(e.status) > 0 {
			if a.status != e.status {
				t.Errorf("%s actual status %s does not equal expected: %s", e.query, a.status, e.status)
			}
		} else if a.term != e.term {
			t.Errorf("%s actual term %s does not equal expected: %s", e.query, a.term, e.term)
		}
	}
}
