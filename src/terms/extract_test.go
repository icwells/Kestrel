// Tests functions in the extract script as well as percent encoding functions

package terms

import (
	"github.com/icwells/go-tools/strarray"
	"github.com/icwells/kestrel/src/kestrelutils"
	"github.com/trustmaster/go-aspell"
	"testing"
)

type extractinput struct {
	corrected string
	query     string
	status    string
	term      string
}

func extractentry(s []string) extractinput {
	// Returns struct from slice input
	var ret extractinput
	ret.query = s[0]
	ret.term = s[1]
	ret.corrected = s[2]
	ret.status = s[3]
	return ret
}

func newExtractInput() []extractinput {
	// Returns initialized slice of test data
	var ret []extractinput
	ret = append(ret, extractentry([]string{"FISH (CARDINAL TETRA OR PENCIL FISH)", "Fish", "", ""}))
	ret = append(ret, extractentry([]string{"PIPING` GUAN ", "Piping guan", "Piping guano", ""}))
	ret = append(ret, extractentry([]string{`SEBA'S  STRIPED FINGERFISH "Sheila"`, "Seba's striped fingerfish", "Sheba's striped finger fish", ""}))
	ret = append(ret, extractentry([]string{"axolotl-5", "Axolotl", "", ""}))
	ret = append(ret, extractentry([]string{"unknown fish", "", "", "uncertainEntry"}))
	ret = append(ret, extractentry([]string{"ferret?", "", "", "uncertainEntry"}))
	ret = append(ret, extractentry([]string{"canine mix", "", "", "hybrid"}))
	ret = append(ret, extractentry([]string{"corgi x", "", "", "hybrid"}))
	ret = append(ret, extractentry([]string{"xy", "", "", "tooShort"}))
	return ret
}

func TestFilter(t *testing.T) {
	speller, _ := aspell.NewSpeller(map[string]string{"lang": "en_US"})
	exp := newExtractInput()
	for _, e := range exp {
		a := NewTerm(e.query)
		a.filter(speller)
		if len(e.status) > 0 {
			if a.Status != e.status {
				t.Errorf("%s actual status %s does not equal expected: %s", e.query, a.Status, e.status)
			}
		} else if a.Term != e.term {
			t.Errorf("%s actual term %s does not equal expected: %s", e.query, a.Term, e.term)
		} else if a.Corrected != e.corrected {
			t.Errorf("%s actual spell-checked term %s does not equal expected: %s", e.query, a.Corrected, e.corrected)
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
		a := strarray.TitleCase(i.input)
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
		encoded := kestrelutils.PercentEncode(e.input)
		if encoded != e.output {
			t.Errorf("Actual encoded string %s does not equal expected: %s", encoded, e.output)
		}
		decoded := kestrelutils.PercentDecode(e.output)
		if decoded != e.input {
			t.Errorf("Actual decoded string %s does not equal expected: %s", decoded, e.input)
		}
	}
}
