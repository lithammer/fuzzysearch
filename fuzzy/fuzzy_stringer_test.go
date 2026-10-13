package fuzzy

import (
	"fmt"
	"sort"
	"strings"
	"testing"
)

type stringer struct {
	string
}

func (s stringer) String() string {
	return s.string
}

var fuzzyTestsStringer = []struct {
	source string
	target fmt.Stringer
	wanted bool
	rank   int
}{
	{"zazz", stringer{deBelloGallico + " zazz"}, true, 1544},
	{"zazz", stringer{"zazz " + deBelloGallico}, true, 1544},
	{"twl", stringer{"cartwheel"}, true, 6},
	{"cart", stringer{"cartwheel"}, true, 5},
	{"cw", stringer{"cartwheel"}, true, 7},
	{"ee", stringer{"cartwheel"}, true, 7},
	{"art", stringer{"cartwheel"}, true, 6},
	{"eeel", stringer{"cartwheel"}, false, -1},
	{"dog", stringer{"cartwheel"}, false, -1},
	{"ёлка", stringer{"ёлочка"}, true, 2},
	{"ветер", stringer{"ёлочка"}, false, -1},
	{"中国", stringer{"中华人民共和国"}, true, 5},
	{"日本", stringer{"中华人民共和国"}, false, -1},
	{"イ", stringer{"イカ"}, true, 1},
	{"limón", stringer{"limon"}, false, -1},
	{"kitten", stringer{"setting"}, false, -1},
}

func TestFuzzyMatchStringer(t *testing.T) {
	for _, val := range fuzzyTestsStringer {
		match := MatchStringer(val.source, val.target)
		if match != val.wanted {
			t.Errorf("%s in %s expected matchStringer to be %t, got %t",
				val.source, val.target, val.wanted, match)
		}
	}
}

func TestFuzzyMatchFoldStringer(t *testing.T) {
	for _, val := range fuzzyTestsStringer {
		match := MatchFoldStringer(val.source, stringer{strings.ToUpper(val.target.String())})
		if match != val.wanted {
			t.Errorf("%s in %s expected matchStringer to be %t, got %t",
				val.source, strings.ToUpper(val.target.String()), val.wanted, match)
		}
	}
}

func TestFuzzyMatchNormalizedStringer(t *testing.T) {
	var normalizedTests = []struct {
		source string
		target fmt.Stringer
		wanted bool
	}{
		{"limon", stringer{"limón"}, true},
		{"limón", stringer{"limon tart"}, true},
		{"limón", stringer{"LiMóN tArT"}, false},
		{"limón", stringer{"LeMoN tArT"}, false},
	}

	for _, val := range normalizedTests {
		match := MatchNormalizedStringer(val.source, val.target)
		if match != val.wanted {
			t.Errorf("%s in %s expected matchStringer to be %t, got %t",
				val.source, val.target, val.wanted, match)
		}
	}
}

func TestFuzzyMatchNormalizedFoldStringer(t *testing.T) {
	var normalizedTests = []struct {
		source string
		target fmt.Stringer
		wanted bool
	}{
		{"limon", stringer{"limón"}, true},
		{"limón", stringer{"limon tart"}, true},
		{"limón", stringer{"LiMóN tArT"}, true},
		{"limón", stringer{"LeMoN tArT"}, false},
	}

	for _, val := range normalizedTests {
		match := MatchNormalizedFoldStringer(val.source, val.target)
		if match != val.wanted {
			t.Errorf("%s in %s expected matchStringer to be %t, got %t",
				val.source, val.target, val.wanted, match)
		}
	}
}

func TestFuzzyFindStringer(t *testing.T) {
	target := []fmt.Stringer{stringer{"cartwheel"}, stringer{"foobar"}, stringer{"wheel"}, stringer{"baz"}, stringer{"cartwhéél"}}
	wanted := []string{"cartwheel", "wheel"}

	matches := FindStringer("whel", target)

	if len(matches) != len(wanted) {
		t.Errorf("expected %s, got %s", wanted, matches)
	}

	for i := range wanted {
		if wanted[i] != matches[i].String() {
			t.Errorf("expected %s, got %s", wanted, matches)
		}
	}
}

func TestFuzzyFindNormalizedStringer(t *testing.T) {
	target := []fmt.Stringer{stringer{"cartwheel"}, stringer{"foobar"}, stringer{"wheel"}, stringer{"baz"}, stringer{"cartwhéél"}, stringer{"WHEEL"}}
	wanted := []string{"cartwheel", "wheel", "cartwhéél"}

	matches := FindNormalizedStringer("whél", target)

	if len(matches) != len(wanted) {
		t.Errorf("expected %s, got %s", wanted, matches)
	}

	for i := range wanted {
		if wanted[i] != matches[i].String() {
			t.Errorf("expected %s, got %s", wanted, matches)
		}
	}
}

func TestFuzzyFindNormalizedFoldStringer(t *testing.T) {
	target := []fmt.Stringer{stringer{"cartwheel"}, stringer{"foobar"}, stringer{"wheel"}, stringer{"baz"}, stringer{"cartwhéél"}, stringer{"WHEEL"}}
	wanted := []string{"cartwheel", "wheel", "cartwhéél", "WHEEL"}

	matches := FindNormalizedFoldStringer("whél", target)

	if len(matches) != len(wanted) {
		t.Errorf("expected %s, got %s", wanted, matches)
	}

	for i := range wanted {
		if wanted[i] != matches[i].String() {
			t.Errorf("expected %s, got %s", wanted, matches)
		}
	}
}

func TestRankMatchStringer(t *testing.T) {
	for _, val := range fuzzyTestsStringer {
		rank := RankMatchStringer(val.source, val.target)
		if rank != val.rank {
			t.Errorf("expected ranking %d, got %d for %s in %s",
				val.rank, rank, val.source, val.target)
		}
	}
}

func TestRankMatchNormalizedStringer(t *testing.T) {
	var fuzzyTests = []struct {
		source string
		target fmt.Stringer
		rank   int
	}{
		{"limó", stringer{"limon"}, 1},
		{"limó", stringer{"limon"}, 1},
		{"limó", stringer{"LIMON"}, -1},
	}

	for _, val := range fuzzyTests {
		rank := RankMatchNormalizedStringer(val.source, val.target)
		if rank != val.rank {
			t.Errorf("expected ranking %d, got %d for %s in %s",
				val.rank, rank, val.source, val.target)
		}
	}
}

func TestRankMatchNormalizedFoldStringer(t *testing.T) {
	var fuzzyTests = []struct {
		source string
		target fmt.Stringer
		rank   int
	}{
		{"limó", stringer{"limon"}, 1},
		{"limó", stringer{"limon"}, 1},
		{"limó", stringer{"LIMON"}, 1},
		{"limó", stringer{"LIMON TART"}, 6},
	}

	for _, val := range fuzzyTests {
		rank := RankMatchNormalizedFoldStringer(val.source, val.target)
		if rank != val.rank {
			t.Errorf("expected ranking %d, got %d for %s in %s",
				val.rank, rank, val.source, val.target)
		}
	}
}

func TestRankMatchNormalizedFoldStringerConcurrent(t *testing.T) {
	var target []fmt.Stringer
	for _, s := range strings.Split("Lorem ipsum dolor sit amet, consectetur adipiscing elit", " ") {
		target = append(target, stringer{s})
	}
	source := "ips"
	procs := 10
	iter := 10
	type empty struct{}
	done := make(chan empty)
	for i := 0; i <= procs; i++ {
		go func() {
			for n := 0; n < iter; n++ {
				_ = RankFindNormalizedFoldStringer(source, target)
			}
			done <- empty{}
		}()
	}
	cnt := 0
	for i := 0; i < procs; i++ {
		<-done
		cnt++
	}
}

func TestRankFindStringer(t *testing.T) {
	target := []fmt.Stringer{stringer{"cartwheel"}, stringer{"foobar"}, stringer{"wheel"}, stringer{"baz"}}
	wanted := []RankStringer{
		{"whl", stringer{"cartwheel"}, 6, 0},
		{"whl", stringer{"wheel"}, 2, 2},
	}

	ranks := RankFindStringer("whl", target)

	if len(ranks) != len(wanted) {
		t.Errorf("expected %+v, got %+v", wanted, ranks)
	}

	for i := range wanted {
		if wanted[i] != ranks[i] {
			t.Errorf("expected %+v, got %+v", wanted, ranks)
		}
	}
}

func TestRankFindNormalizedStringer(t *testing.T) {
	target := []fmt.Stringer{stringer{"limón"}, stringer{"limon"}, stringer{"lemon"}, stringer{"LIMON"}}
	wanted := []RankStringer{
		{"limó", stringer{"limón"}, 1, 0},
		{"limó", stringer{"limon"}, 2, 1},
	}

	ranks := RankFindNormalizedStringer("limó", target)

	if len(ranks) != len(wanted) {
		t.Errorf("expected %+v, got %+v", wanted, ranks)
	}

	for i := range wanted {
		if wanted[i] != ranks[i] {
			t.Errorf("expected %+v, got %+v", wanted, ranks)
		}
	}
}

func TestRankFindNormalizedFoldStringer(t *testing.T) {
	target := []fmt.Stringer{stringer{"limón"}, stringer{"limon"}, stringer{"lemon"}, stringer{"LIMON"}}
	wanted := []RankStringer{
		{"limó", stringer{"limón"}, 1, 0},
		{"limó", stringer{"limon"}, 2, 1},
		{"limó", stringer{"LIMON"}, 5, 3},
	}

	ranks := RankFindNormalizedFoldStringer("limó", target)

	if len(ranks) != len(wanted) {
		t.Errorf("expected %+v, got %+v", wanted, ranks)
	}

	for i := range wanted {
		if wanted[i] != ranks[i] {
			t.Errorf("expected %+v, got %+v", wanted, ranks)
		}
	}
}

func TestSortingRanksStringer(t *testing.T) {
	rs := RanksStringer{{"a", stringer{"b"}, 1, 0}, {"a", stringer{"cc"}, 2, 1}, {"a", stringer{"a"}, 0, 2}}
	wanted := RanksStringer{rs[2], rs[0], rs[1]}

	sort.Sort(rs)

	for i := range wanted {
		if wanted[i] != rs[i] {
			t.Errorf("expected %+v, got %+v", wanted, rs)
		}
	}
}

func BenchmarkMatchStringer(b *testing.B) {
	ft := fuzzyTestsStringer[2]
	for i := 0; i < b.N; i++ {
		MatchStringer(ft.source, ft.target)
	}
}

func BenchmarkMatchStringerBigLate(b *testing.B) {
	ft := fuzzyTestsStringer[0]
	for i := 0; i < b.N; i++ {
		MatchStringer(ft.source, ft.target)
	}
}

func BenchmarkMatchStringerBigEarly(b *testing.B) {
	ft := fuzzyTestsStringer[1]
	for i := 0; i < b.N; i++ {
		MatchStringer(ft.source, ft.target)
	}
}

func BenchmarkMatchFoldStringer(b *testing.B) {
	ft := fuzzyTestsStringer[2]
	for i := 0; i < b.N; i++ {
		MatchFoldStringer(ft.source, ft.target)
	}
}

func BenchmarkMatchFoldStringerBigLate(b *testing.B) {
	ft := fuzzyTestsStringer[0]
	for i := 0; i < b.N; i++ {
		MatchFoldStringer(ft.source, ft.target)
	}
}

func BenchmarkMatchFoldStringerBigEarly(b *testing.B) {
	ft := fuzzyTestsStringer[1]
	for i := 0; i < b.N; i++ {
		MatchFoldStringer(ft.source, ft.target)
	}
}

func BenchmarkRankMatchStringer(b *testing.B) {
	ft := fuzzyTestsStringer[2]
	for i := 0; i < b.N; i++ {
		RankMatchStringer(ft.source, ft.target)
	}
}

func BenchmarkRankMatchStringerBigLate(b *testing.B) {
	ft := fuzzyTestsStringer[0]
	for i := 0; i < b.N; i++ {
		RankMatchStringer(ft.source, ft.target)
	}
}

func BenchmarkRankMatchStringerBigEarly(b *testing.B) {
	ft := fuzzyTestsStringer[1]
	for i := 0; i < b.N; i++ {
		RankMatchStringer(ft.source, ft.target)
	}
}

func ExampleMatchStringer() {
	fmt.Print(MatchStringer("twl", stringer{
		"cartwheel",
	}))
	// Output: true
}

func ExampleFindStringer() {
	fmt.Print(FindStringer("whl", []fmt.Stringer{stringer{"cartwheel"}, stringer{"foobar"}, stringer{"wheel"}, stringer{"baz"}}))
	// Output: [cartwheel wheel]
}

func ExampleRankMatchStringer() {
	fmt.Print(RankMatchStringer("twl",
		stringer{
			"cartwheel",
		}))
	// Output: 6
}

func ExampleRankFindStringer() {
	fmt.Printf("%+v", RankFindStringer("whl", []fmt.Stringer{stringer{"cartwheel"}, stringer{"foobar"}, stringer{"wheel"}, stringer{"baz"}}))
	// Output: [{Source:whl Target:cartwheel Distance:6 OriginalIndex:0} {Source:whl Target:wheel Distance:2 OriginalIndex:2}]
}
