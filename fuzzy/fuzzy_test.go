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
	{"stakes", "The only really major change from the <em>status quo ante</em> resulting from the War of Austrian Succession was Austria’s loss of Silesia to Prussia. As a consequence, Frederick increased the number of his subjects by over 50 percent and made Prussia the first new entrant to the ranks of the Great Powers since the Peace of Westphalia in 1648. Frederick—who by this point was already being called “the Great” by his admirers—had demonstrated that as a political instrument war could accomplish major objectives even in an age of “limited war,” when sociopolitical constraints as well as military ones tended to make wars relatively indecisive and profitless. Still, the Prussian king was able to make such substantial gains only because his superior soldiers and his increasingly adept generalship enabled him to win five consecutive battles against the Austrians. Even that impressive military record, however, bought Prussia only an insecure conquest and a cold and temporary peace. Maria Theresa remained determined to recover Silesia, and worked tirelessly to improve her army and adjust the political situation to make that possible. The result, the subject of chapter 9, was a renewed effort to undo Frederick’s gains, a war that would push Austria, France, Britain, and Prussia towards a style of warfare employing greater resources and involving higher political stakes.", true, 1374},
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
	ft := fuzzyTests[0]
	for i := 0; i < b.N; i++ {
		Match(ft.source, ft.target)
	}
}

func BenchmarkRankMatch(b *testing.B) {
	ft := fuzzyTests[0]
	for i := 0; i < b.N; i++ {
		RankMatch(ft.source, ft.target)
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
