// Tests scorer struct and methods

package searchtaxa

import (
	"testing"
)

func TestGetMax(t *testing.T) {
	s := newScorer()
	s.scores["a"] = make(map[string]int)
	s.scores["b"] = make(map[string]int)
	s.scores["a"]["b"] = 5
	s.scores["a"]["c"] = 7
	s.scores["b"]["c"] = 6
	a1, a2 := s.getMax()
	if a1 != "a" || a2 != "c" {
		t.Errorf("Actual match keys %s %s do not equal expected: ac", a1, a2)
	}
}

func TestScoreLevel(t *testing.T) {
	s := newScorer()
	levels := []struct {
		t1, t2   string
		expected int
	}{
		{"Canis", "Canis", 1},
		{"Squamata", "Reptilia", -1},
		{"Canis", "NA", 0},
	}
	for _, i := range levels {
		a := s.scoreLevel(i.t1, i.t2)
		if a != i.expected {
			t.Errorf("Actual score for %s and %s %d does not equal expected %d", i.t1, i.t2, a, i.expected)
		}
	}
}
