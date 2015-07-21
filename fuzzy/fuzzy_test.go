package fuzzy

import (
	"fmt"
	"testing"
)

var fuzzyTests = []struct {
	needle   string
	haystack string
	wanted   bool
}{
	{"twl", "cartwheel", true},
	{"cart", "cartwheel", true},
	{"cw", "cartwheel", true},
	{"ee", "cartwheel", true},
	{"art", "cartwheel", true},
	{"eeel", "cartwheel", false},
	{"dog", "cartwheel", false},
}

func TestFuzzyMatch(t *testing.T) {
	for _, v := range fuzzyTests {
		match := Search(v.needle, v.haystack)
		if match != v.wanted {
			t.Errorf("%s in %s expected match to be %t, got %t", v.needle, v.haystack, v.wanted, match)
		}
	}
}

func BenchmarkFuzzySearch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Search("twl", "cartwheel")
	}
}

func ExampleFuzzySearch() {
	fmt.Print(Search("twl", "cartwheel"))
	// Output: true
}
