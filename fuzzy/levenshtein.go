package fuzzy

// LevenshteinDistance measures the difference between two strings.
// The Levenshtein distance between two words is the minimum number of
// single-character edits (i.e. insertions, deletions or substitutions)
// required to change one word into the other.
//
// This implemention is optimized to use O(min(m,n)) space and is based on the
// optimized C version found here:
// http://en.wikibooks.org/wiki/Algorithm_implementation/Strings/Levenshtein_distance#C
func LevenshteinDistance(s1, s2 string) int {
	column := make([]int, len(s1)+1)

	for y := 1; y <= len(s1); y++ {
		column[y] = y
	}

	for x := 1; x <= len(s2); x++ {
		column[0] = x

		for y, lastDiag := 1, x-1; y <= len(s1); y++ {
			oldDiag := column[y]
			cost := 0
			if s1[y-1] != s2[x-1] {
				cost = 1
			}
			column[y] = min(column[y]+1, column[y-1]+1, lastDiag+cost)
			lastDiag = oldDiag
		}
	}

	return column[len(s1)]
}

func min(a, b, c int) int {
	if a < b && a < c {
		return a
	} else if b < c {
		return b
	}
	return c
}
