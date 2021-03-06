// Tests functions in the extract script as well as percent encoding functions

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
	ret = append(ret, extractentry([]string{"axolotl-5", "Axolotl", ""}))
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

func TestTitleCase(t *testing.T) {
	str := []struct {
		input, expected string
	}{
		{"SEBA'S STRIPED  FINGERFISH", "Seba's Striped Fingerfish"},
		{"Sharp shinned Hawk", "Sharp Shinned Hawk"},
		{"PIPING` x GUAN ", "Piping` Guan"},
	}
	for _, i := range str {
		a := titleCase(i.input)
		if a != i.expected {
			t.Errorf("Actual term %s does not equal expected: %s", a, i.expected)
		}
	}
}

func TestEncoding(t *testing.T) {
	// Tests percent encode and decode functions
	expected := []struct {
		input, output string
	}{
		{"SEBA'S STRIPED FINGERFISH", "SEBA%27S%20STRIPED%20FINGERFISH"},
		{"Sharp shinned Hawk", "Sharp%20shinned%20Hawk"},
	}
	for _, e := range expected {
		encoded := percentEncode(e.input)
		if encoded != e.output {
			t.Errorf("Actual encoded string %s does not equal expected: %s", encoded, e.output)
		}
		decoded := percentDecode(e.output)
		if decoded != e.input {
			t.Errorf("Actual decoded string %s does not equal expected: %s", decoded, e.input)
		}
	}
}
