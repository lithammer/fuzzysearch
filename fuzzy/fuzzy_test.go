package fuzzy

import (
	"fmt"
	"sort"
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
	{"ёлка", "ёлочка", true, 2},
	{"ветер", "ёлочка", false, -1},
	{"中国", "中华人民共和国", true, 5},
	{"日本", "中华人民共和国", false, -1},
}

func TestFuzzyMatch(t *testing.T) {
	for _, val := range fuzzyTests {
		match := Match(val.needle, val.haystack)
		if match != val.wanted {
			t.Errorf("%s in %s expected match to be %t, got %t",
				val.needle, val.haystack, val.wanted, match)
		}
	}
}

func TestFuzzyFind(t *testing.T) {
	haystack := []string{"cartwheel", "foobar", "wheel", "baz"}
	wanted := []string{"cartwheel", "wheel"}

	matches := Find("whl", haystack)

	if len(matches) != len(wanted) {
		t.Errorf("expected %s, got %s", wanted, matches)
	}

	for i := range wanted {
		if wanted[i] != matches[i] {
			t.Errorf("expected %s, got %s", wanted, matches)
		}
	}
}

func TestRankMatch(t *testing.T) {
	for _, val := range fuzzyTests {
		rank := RankMatch(val.needle, val.haystack)
		if rank != val.rank {
			t.Errorf("expected ranking %d, got %d for %s in %s", val.rank, rank, val.needle, val.haystack)
		}
	}
}

func TestRankFind(t *testing.T) {
	haystack := []string{"cartwheel", "foobar", "wheel", "baz"}
	wanted := []Rank{
		{"cartwheel", 6},
		{"wheel", 2},
	}

	ranks := RankFind("whl", haystack)

	if len(ranks) != len(wanted) {
		t.Errorf("expected %+v, got %+v", wanted, ranks)
	}

	for i := range wanted {
		if wanted[i] != ranks[i] {
			t.Errorf("expected %+v, got %+v", wanted, ranks)
		}
	}
}

func TestSortingRanks(t *testing.T) {
	ranks := Ranks{{"b", 3}, {"c", 6}, {"a", 2}}
	wanted := Ranks{{"a", 2}, {"b", 3}, {"c", 6}}

	sort.Sort(ranks)

	for i := range wanted {
		if wanted[i] != ranks[i] {
			t.Errorf("expected %+v, got %+v", wanted, ranks)
		}
	}
}

func BenchmarkMatch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Match("kitten", "sitting")
	}
}

func BenchmarkRankMatch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RankMatch("kitten", "sitting")
	}
}

func ExampleMatch() {
	fmt.Print(Match("twl", "cartwheel"))
	// Output: true
}

func ExampleFind() {
	fmt.Print(Find("whl", []string{"cartwheel", "foobar", "wheel", "baz"}))
	// Output: [cartwheel wheel]
}

func ExampleRankMatch() {
	fmt.Print(RankMatch("twl", "cartwheel"))
	// Output: 6
}

func ExampleRankFind() {
	fmt.Print(RankFind("whl", []string{"cartwheel", "foobar", "wheel", "baz"}))
	// Output: [{cartwheel 6} {wheel 2}]
}
