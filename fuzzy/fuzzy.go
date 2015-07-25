package fuzzy

// Match does a partial match ("fuzzy search") of needle in haystack.
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

type Ranks []Rank

func (r Ranks) Len() int {
	return len(r)
}

func (r Ranks) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r Ranks) Less(i, j int) bool {
	return r[i].Distance < r[j].Distance
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
func RankFind(needle string, haystacks []string) Ranks {
	var ranks Ranks
	for _, val := range Find(needle, haystacks) {
		ranks = append(ranks, Rank{val, LevenshteinDistance(needle, val)})
	}
	return ranks
}
