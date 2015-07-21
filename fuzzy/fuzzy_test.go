package fuzzy

import (
	"fmt"
	"testing"
)

type testVector struct {
	needle   string
	haystack string
	match    bool
}

var fuzzyTestVectors = []testVector{
	{"twl", "cartwheel", true},
	{"cart", "cartwheel", true},
	{"cw", "cartwheel", true},
	{"ee", "cartwheel", true},
	{"art", "cartwheel", true},
	{"eeel", "cartwheel", false},
	{"dog", "cartwheel", false},
}

func TestFuzzyMatch(t *testing.T) {
	for _, v := range fuzzyTestVectors {
		m := Search(v.needle, v.haystack)
		if m != v.match {
			t.Errorf("%s in %s expected match to be %t, got %t", v.needle, v.haystack, v.match, m)
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
