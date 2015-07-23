package fuzzy

// Search does a partial match ("fuzzy search") of needle in haystack.
func Search(needle, haystack string) bool {
	nlen := len(needle)
	hlen := len(haystack)

	if nlen > hlen {
		return false
	}

	if nlen == hlen {
		return needle == haystack
	}

Outer:
	for i, j := 0, 0; i < nlen; i++ {
		for j < hlen {
			j++
			if needle[i] == haystack[j-1] {
				continue Outer
			}
		}
		return false
	}

	return true
}

// SearchMany will return a list of strings in haysatcks that fuzzy matches
// needle.
func SearchMany(needle string, haystacks []string) []string {
	var matches []string

	for _, haystack := range haystacks {
		if Search(needle, haystack) {
			matches = append(matches, haystack)
		}
	}

	return matches
}

type Rank struct {
	Word     string
	Distance int
}

// RankSearch is similar to Search except it will measure the Levenshtein
// distance between the needle and the haystack and return its result.
// If there was no match, it will return -1.
func RankSearch(needle, haystack string) int {
	match := Search(needle, haystack)
	if !match {
		return -1
	}
	return LevenshteinDistance(needle, haystack)
}

// RankSearchMany is similar to SearchMany, except it will also rank all matches
// using Levenshtein distance.
func RankSearchMany(needle string, haystacks []string) []Rank {
	var ranks []Rank
	for _, val := range SearchMany(needle, haystacks) {
		ranks = append(ranks, Rank{val, LevenshteinDistance(needle, val)})
	}
	return ranks
}
