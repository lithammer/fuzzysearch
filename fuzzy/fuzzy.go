package fuzzy

// Search does a partial match ("fuzzy search") of needle in haystack
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
