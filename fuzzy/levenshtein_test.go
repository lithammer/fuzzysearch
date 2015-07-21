package fuzzy

import "testing"

var levenshteinDistanceTests = []struct {
	s1     string
	s2     string
	wanted int
}{
	{"a", "a", 0},
	{"ab", "ab", 0},
	{"ab", "aa", 1},
	{"ab", "aa", 1},
	{"ab", "aaa", 2},
	{"bbb", "a", 3},
	{"kitten", "sitting", 3},
}

func TestLevenshtein(t *testing.T) {
	for _, test := range levenshteinDistanceTests {
		distance := LevenshteinDistance(test.s1, test.s2)
		if distance != test.wanted {
			t.Errorf("got distance %d, expected %d", distance, test.wanted)
		}
	}
}

func BenchmarkLevenshteinDistance(b *testing.B) {
	for i := 0; i < b.N; i++ {
		LevenshteinDistance("aaa", "aba")
		LevenshteinDistance("kitten", "sitting")
	}
}
