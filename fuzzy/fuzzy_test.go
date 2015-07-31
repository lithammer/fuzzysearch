package fuzzy

import (
	"fmt"
	"sort"
	"testing"
)

var fuzzyTests = []struct {
	source string
	target string
	wanted bool
	rank   int
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
		match := Match(val.source, val.target)
		if match != val.wanted {
			t.Errorf("%s in %s expected match to be %t, got %t",
				val.source, val.target, val.wanted, match)
		}
	}
}

func TestFuzzyFind(t *testing.T) {
	target := []string{"cartwheel", "foobar", "wheel", "baz"}
	wanted := []string{"cartwheel", "wheel"}

	matches := Find("whl", target)

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
		rank := RankMatch(val.source, val.target)
		if rank != val.rank {
			t.Errorf("expected ranking %d, got %d for %s in %s", val.rank, rank, val.source, val.target)
		}
	}
}

func TestRankFind(t *testing.T) {
	target := []string{"cartwheel", "foobar", "wheel", "baz"}
	wanted := []Rank{
		{"whl", "cartwheel", 6},
		{"whl", "wheel", 2},
	}

	ranks := RankFind("whl", target)

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
	rs := ranks{{"a", "b", 1}, {"a", "cc", 2}, {"a", "a", 0}}
	wanted := ranks{rs[2], rs[0], rs[1]}

	sort.Sort(rs)

	for i := range wanted {
		if wanted[i] != rs[i] {
			t.Errorf("expected %+v, got %+v", wanted, rs)
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
	fmt.Printf("%+v", RankFind("whl", []string{"cartwheel", "foobar", "wheel", "baz"}))
	// Output: [{Source:whl Target:cartwheel Distance:6} {Source:whl Target:wheel Distance:2}]
}
