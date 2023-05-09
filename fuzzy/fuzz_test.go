package fuzzy

import (
	"sort"
	"testing"
)

func FuzzFind(f *testing.F) {
	f.Fuzz(func(t *testing.T, n string, h []byte) {
		s := make([]string, len(h))
		for i, b := range h {
			s[i] = string(b)
		}
		Find(n, s)
		FindFold(n, s)
		FindNormalized(n, s)
		FindNormalizedFold(n, s)
		r := RankFind(n, s)
		sort.Sort(r)
		// No need to sort the other Rank calls;
		// assume first sort can catch any bugs.
		RankFindFold(n, s)
		RankFindNormalized(n, s)
		RankFindNormalizedFold(n, s)
		if len(s) > 0 {
			x := s[0]
			LevenshteinDistance(n, x)
			Match(n, x)
			MatchFold(n, x)
			MatchNormalized(n, x)
			MatchNormalizedFold(n, x)
		}
	})
}
