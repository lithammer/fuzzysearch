package fuzzy

import "unicode/utf8"

// Match returns true if needle matches haystack using a fuzzy-searching
// algorithm. Note that it doesn't implement Levenshtein distance (see
// RankMatch instead), but rather a simplified version where there's no
// approximation. The method will return true only if each character in the
// needle can be found in the haystack and occurs after the preceding matches.
func Match(needle, haystack string) bool {
	nlen := len(needle)
	hlen := len(haystack)

	if nlen > hlen {
		return false
	}

	if nlen == hlen {
		return needle == haystack
	}
Outer:
	for _, r1 := range needle {
		for i, r2 := range haystack {
			if r1 == r2 {
				haystack = haystack[i+utf8.RuneLen(r2):]
				continue Outer
			}
		}
		return false
	}

	return true
}

// Find will return a list of strings in haysatcks that fuzzy matches
// needle.
func Find(needle string, haystacks []string) []string {
	var matches []string

	for _, haystack := range haystacks {
		if Match(needle, haystack) {
			matches = append(matches, haystack)
		}
	}

	return matches
}

type Rank struct {
	Word     string
	Distance int
}

// RankMatch is similar to Match except it will measure the Levenshtein
// distance between the needle and the haystack and return its result.
// If there was no match, it will return -1.
func RankMatch(needle, haystack string) int {
	match := Match(needle, haystack)
	if !match {
		return -1
	}
	return LevenshteinDistance(needle, haystack)
}

// RankFind is similar to Find, except it will also rank all matches
// using Levenshtein distance.
func RankFind(needle string, haystacks []string) ranks {
	var r ranks
	for _, word := range Find(needle, haystacks) {
		r = append(r, Rank{
			Word:   word,
			Distance: LevenshteinDistance(needle, word),
		})
	}
	return r
}

type ranks []Rank

func (r ranks) Len() int {
	return len(r)
}

func (r ranks) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r ranks) Less(i, j int) bool {
	return r[i].Distance < r[j].Distance
}
