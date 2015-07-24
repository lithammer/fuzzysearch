package fuzzy

import (
	"fmt"
	"testing"
)

var fuzzyTests = []struct {
	needle   string
	haystack string
	wanted   bool
	rank     int
}{
	{"twl", "cartwheel", true, 6},
	{"cart", "cartwheel", true, 5},
	{"cw", "cartwheel", true, 7},
	{"ee", "cartwheel", true, 7},
	{"art", "cartwheel", true, 6},
	{"eeel", "cartwheel", false, -1},
	{"dog", "cartwheel", false, -1},
}

func TestFuzzySearch(t *testing.T) {
	for _, val := range fuzzyTests {
		match := Search(val.needle, val.haystack)
		if match != val.wanted {
			t.Errorf("%s in %s expected match to be %t, got %t",
				val.needle, val.haystack, val.wanted, match)
		}
	}
}

func TestFuzzySearchMany(t *testing.T) {
	haystack := []string{"cartwheel", "foobar", "wheel", "baz"}
	wanted := []string{"cartwheel", "wheel"}

	matches := SearchMany("whl", haystack)

	if len(matches) != len(wanted) {
		t.Errorf("expected %s, got %s", wanted, matches)
	}

	for i := range wanted {
		if wanted[i] != matches[i] {
			t.Errorf("expected %s, got %s", wanted, matches)
		}
	}
}

func TestRankSearch(t *testing.T) {
	for _, val := range fuzzyTests {
		rank := RankSearch(val.needle, val.haystack)
		if rank != val.rank {
			t.Errorf("expected ranking %d, got %d", val.rank, rank)
		}
	}
}

func TestRankSearchMany(t *testing.T) {
	haystack := []string{"cartwheel", "foobar", "wheel", "baz"}
	wanted := []Rank{
		{"cartwheel", 6},
		{"wheel", 2},
	}

	ranks := RankSearchMany("whl", haystack)

	if len(ranks) != len(wanted) {
		t.Errorf("expected %+v, got %+v", wanted, ranks)
	}

	for i := range wanted {
		if wanted[i] != ranks[i] {
			t.Errorf("expected %+v, got %+v", wanted, ranks)
		}
	}
}

func BenchmarkFuzzySearch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Search("twl", "cartwheel")
	}
}

func ExampleSearch() {
	fmt.Print(Search("twl", "cartwheel"))
	// Output: true
}

func ExampleSearchMany() {
	fmt.Print(SearchMany("whl", []string{"cartwheel", "foobar", "wheel", "baz"}))
	// Output: [cartwheel wheel]
}

func ExampleRankSearch() {
	fmt.Print(RankSearch("twl", "cartwheel"))
	// Output: 6
}

func ExampleRankSearchMany() {
	fmt.Print(RankSearchMany("whl", []string{"cartwheel", "foobar", "wheel", "baz"}))
	// Output: [{cartwheel 6} {wheel 2}]
}
