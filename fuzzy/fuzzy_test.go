package fuzzy

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
)

const deBelloGallico = `All Gaul is divided into three parts, one of which the Belgae inhabit,
the Aquitani another, those who in their own language are called Celts, in our Gauls, the third.
All these differ from each other in language, customs and laws. The river Garonne separates the
Gauls from the Aquitani; the Marne and the Seine separate them from the Belgae. Of all these,
the Belgae are the bravest, because they are furthest from the civilization and refinement of
[our] Province, and merchants least frequently resort to them, and import those things which tend
to effeminate the mind; and they are the nearest to the Germans, who dwell beyond the Rhine,
with whom they are continually waging war; for which reason the Helvetii also surpass the rest
of the Gauls in valor, as they contend with the Germans in almost daily battles, when they either
repel them from their own territories, or themselves wage war on their frontiers. One part of
these, which it has been said that the Gauls occupy, takes its beginning at the river Rhone;
it is bounded by the river Garonne, the ocean, and the territories of the Belgae; it borders,
too, on the side of the Sequani and the Helvetii, upon the river Rhine, and stretches toward the
north. The Belgae rises from the extreme frontier of Gaul, extend to the lower part of the river
Rhine; and look toward the north and the rising sun. Aquitania extends from the river Garonne to
the Pyrenaean mountains and to that part of the ocean which is near Spain: it looks between the
setting of the sun, and the north star.`

type fuzzyTest struct {
	source string
	target string
	wanted bool
	rank   int
}

var fuzzyTests = []fuzzyTest{
	{"zazz", deBelloGallico + " zazz", true, 1544},
	{"zazz", "zazz " + deBelloGallico, true, 1544},
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
	{"イ", "イカ", true, 1},
	{"limón", "limon", false, -1},
	{"kitten", "setting", false, -1},
	{"\xffinvalid UTF-8\xff", "", false, -1}, // invalid UTF-8
	{"Ⱦ", "", false, -1},                     // uppercase and lowercase runes have different UTF-8 encoding lengths
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

func TestFuzzyMatchFold(t *testing.T) {
	for _, val := range fuzzyTests {
		match := MatchFold(val.source, strings.ToUpper(val.target))
		if match != val.wanted {
			t.Errorf("%s in %s expected match to be %t, got %t",
				val.source, strings.ToUpper(val.target), val.wanted, match)
		}
	}
}

func TestFuzzyMatchNormalized(t *testing.T) {
	normalizedTests := []struct {
		source string
		target string
		wanted bool
	}{
		{"limon", "limón", true},
		{"limón", "limon tart", true},
		{"limón", "LiMóN tArT", false},
		{"limón", "LeMoN tArT", false},
	}

	for _, val := range normalizedTests {
		match := MatchNormalized(val.source, val.target)
		if match != val.wanted {
			t.Errorf("%s in %s expected match to be %t, got %t",
				val.source, val.target, val.wanted, match)
		}
	}
}

func TestFuzzyMatchNormalizedFold(t *testing.T) {
	normalizedTests := []struct {
		source string
		target string
		wanted bool
	}{
		{"limon", "limón", true},
		{"limón", "limon tart", true},
		{"limón", "LiMóN tArT", true},
		{"limón", "LeMoN tArT", false},
	}

	for _, val := range normalizedTests {
		match := MatchNormalizedFold(val.source, val.target)
		if match != val.wanted {
			t.Errorf("%s in %s expected match to be %t, got %t",
				val.source, val.target, val.wanted, match)
		}
	}
}

func TestMatchIndices(t *testing.T) {
	tests := []struct {
		source string
		target string
		wanted []int
	}{
		{"twl", "cartwheel", []int{3, 4, 8}},
		{"cart", "cartwheel", []int{0, 1, 2, 3}},
		{"cw", "cartwheel", []int{0, 4}},
		{"ee", "cartwheel", []int{6, 7}},
		{"art", "cartwheel", []int{1, 2, 3}},
		{"zazz", "zazz " + deBelloGallico, []int{0, 1, 2, 3}},
		{"中国", "中华人民共和国", []int{0, 6}},
		{"ёлка", "ёлочка", []int{0, 1, 4, 5}},
		{"イ", "イカ", []int{0}},
		{"eeel", "cartwheel", nil},
		{"dog", "cartwheel", nil},
		{"kitten", "setting", nil},
		{"ветер", "ёлочка", nil},
		{"日本", "中华人民共和国", nil},
		{"", "cartwheel", []int{}},
		{"", "", []int{}},
		{"a", "", nil},
		{"cartwheel", "", nil},
		{"cartwheels", "cartwheel", nil},
		{"\xffinvalid UTF-8\xff", "", nil},
		{"Ⱦ", "", nil},
	}

	for _, tt := range tests {
		got := MatchIndices(tt.source, tt.target)
		if !reflect.DeepEqual(got, tt.wanted) {
			t.Errorf("MatchIndices(%q, %q) = %v; want %v", tt.source, tt.target, got, tt.wanted)
		}
	}
}

func TestMatchIndicesFold(t *testing.T) {
	tests := []struct {
		source string
		target string
		wanted []int
	}{
		{"twl", "CARTWHEEL", []int{3, 4, 8}},
		{"TWL", "cartwheel", []int{3, 4, 8}},
		{"cw", "CartWheel", []int{0, 4}},
		{"CW", "cartwheel", []int{0, 4}},
		{"ART", "cartwheel", []int{1, 2, 3}},
		{"中国", "中华人民共和国", []int{0, 6}},
		{"dog", "CARTWHEEL", nil},
		{"", "CARTWHEEL", []int{}},
		{"", "", []int{}},
		{"A", "", nil},
	}

	for _, tt := range tests {
		got := MatchIndicesFold(tt.source, tt.target)
		if !reflect.DeepEqual(got, tt.wanted) {
			t.Errorf("MatchIndicesFold(%q, %q) = %v; want %v", tt.source, tt.target, got, tt.wanted)
		}
	}
}

func TestMatchIndicesNormalized(t *testing.T) {
	tests := []struct {
		source string
		target string
		wanted []int
	}{
		{"limon", "limón", []int{0, 1, 2, 3, 4}},
		{"limón", "limon tart", []int{0, 1, 2, 3, 4}},
		{"cartwheel", "cartwhéél", []int{0, 1, 2, 3, 4, 5, 6, 7, 8}},
		{"limón", "LiMóN tArT", nil},
		{"limón", "LeMoN tArT", nil},
		{"", "limón", []int{}},
		{"", "", []int{}},
		{"limón", "", nil},
		{"limon", "limo\u0301n", []int{0, 1, 2, 3, 5}},
	}

	for _, tt := range tests {
		got := MatchIndicesNormalized(tt.source, tt.target)
		if !reflect.DeepEqual(got, tt.wanted) {
			t.Errorf("MatchIndicesNormalized(%q, %q) = %v; want %v", tt.source, tt.target, got, tt.wanted)
		}
	}
}

func TestMatchIndicesNormalizedFold(t *testing.T) {
	tests := []struct {
		source string
		target string
		wanted []int
	}{
		{"limon", "limón", []int{0, 1, 2, 3, 4}},
		{"limón", "limon tart", []int{0, 1, 2, 3, 4}},
		{"limón", "LiMóN tArT", []int{0, 1, 2, 3, 4}},
		{"LIMON", "limón", []int{0, 1, 2, 3, 4}},
		{"limón", "LeMoN tArT", nil},
		{"", "LiMóN", []int{}},
		{"", "", []int{}},
		{"LiMóN", "", nil},
		{"limon", "LIMO\u0301N", []int{0, 1, 2, 3, 5}},
	}

	for _, tt := range tests {
		got := MatchIndicesNormalizedFold(tt.source, tt.target)
		if !reflect.DeepEqual(got, tt.wanted) {
			t.Errorf("MatchIndicesNormalizedFold(%q, %q) = %v; want %v", tt.source, tt.target, got, tt.wanted)
		}
	}
}

func TestMatchIndicesParity(t *testing.T) {
	for _, val := range fuzzyTests {
		match := Match(val.source, val.target)
		indices := MatchIndices(val.source, val.target)
		if (indices != nil) != match {
			t.Errorf("MatchIndices parity mismatch for %q in %q: Match=%t, indices=%v",
				val.source, val.target, match, indices)
		}

		matchFold := MatchFold(val.source, strings.ToUpper(val.target))
		indicesFold := MatchIndicesFold(val.source, strings.ToUpper(val.target))
		if (indicesFold != nil) != matchFold {
			t.Errorf("MatchIndicesFold parity mismatch for %q in %q: MatchFold=%t, indicesFold=%v",
				val.source, strings.ToUpper(val.target), matchFold, indicesFold)
		}
	}
}

func TestFuzzyFind(t *testing.T) {
	target := []string{"cartwheel", "foobar", "wheel", "baz", "cartwhéél"}
	wanted := []string{"cartwheel", "wheel"}

	matches := Find("whel", target)

	if len(matches) != len(wanted) {
		t.Errorf("expected %s, got %s", wanted, matches)
	}

	for i := range wanted {
		if wanted[i] != matches[i] {
			t.Errorf("expected %s, got %s", wanted, matches)
		}
	}
}

func TestFuzzyFindNormalized(t *testing.T) {
	target := []string{"cartwheel", "foobar", "wheel", "baz", "cartwhéél", "WHEEL"}
	wanted := []string{"cartwheel", "wheel", "cartwhéél"}

	matches := FindNormalized("whél", target)

	if len(matches) != len(wanted) {
		t.Errorf("expected %s, got %s", wanted, matches)
	}

	for i := range wanted {
		if wanted[i] != matches[i] {
			t.Errorf("expected %s, got %s", wanted, matches)
		}
	}
}

func TestFuzzyFindNormalizedFold(t *testing.T) {
	target := []string{"cartwheel", "foobar", "wheel", "baz", "cartwhéél", "WHEEL"}
	wanted := []string{"cartwheel", "wheel", "cartwhéél", "WHEEL"}

	matches := FindNormalizedFold("whél", target)

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
			t.Errorf("expected ranking %d, got %d for %s in %s",
				val.rank, rank, val.source, val.target)
		}
	}
}

func TestRankMatchNormalized(t *testing.T) {
	fuzzyTests := []struct {
		source string
		target string
		rank   int
	}{
		{"limó", "limon", 1},
		{"limó", "limon", 1},
		{"limó", "LIMON", -1},
	}

	for _, val := range fuzzyTests {
		rank := RankMatchNormalized(val.source, val.target)
		if rank != val.rank {
			t.Errorf("expected ranking %d, got %d for %s in %s",
				val.rank, rank, val.source, val.target)
		}
	}
}

func TestRankMatchNormalizedFold(t *testing.T) {
	fuzzyTests := []struct {
		source string
		target string
		rank   int
	}{
		{"limó", "limon", 1},
		{"limó", "limon", 1},
		{"limó", "LIMON", 1},
		{"limó", "LIMON TART", 6},
	}

	for _, val := range fuzzyTests {
		rank := RankMatchNormalizedFold(val.source, val.target)
		if rank != val.rank {
			t.Errorf("expected ranking %d, got %d for %s in %s",
				val.rank, rank, val.source, val.target)
		}
	}
}

func TestRankMatchNormalizedFoldConcurrent(t *testing.T) {
	target := strings.Split("Lorem ipsum dolor sit amet, consectetur adipiscing elit", " ")
	source := "ips"
	procs := 10
	iter := 10
	type empty struct{}
	done := make(chan empty)
	for i := 0; i <= procs; i++ {
		go func() {
			for range iter {
				_ = RankFindNormalizedFold(source, target)
			}
			done <- empty{}
		}()
	}
	cnt := 0
	for range procs {
		<-done
		cnt++
	}
}

func TestRankFind(t *testing.T) {
	target := []string{"cartwheel", "foobar", "wheel", "baz"}
	wanted := []Rank{
		{"whl", "cartwheel", 6, 0},
		{"whl", "wheel", 2, 2},
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

func TestRankFindNormalized(t *testing.T) {
	target := []string{"limón", "limon", "lemon", "LIMON"}
	wanted := []Rank{
		{"limó", "limón", 1, 0},
		{"limó", "limon", 2, 1},
	}

	ranks := RankFindNormalized("limó", target)

	if len(ranks) != len(wanted) {
		t.Errorf("expected %+v, got %+v", wanted, ranks)
	}

	for i := range wanted {
		if wanted[i] != ranks[i] {
			t.Errorf("expected %+v, got %+v", wanted, ranks)
		}
	}
}

func TestRankFindNormalizedFold(t *testing.T) {
	target := []string{"limón", "limon", "lemon", "LIMON"}
	wanted := []Rank{
		{"limó", "limón", 1, 0},
		{"limó", "limon", 2, 1},
		{"limó", "LIMON", 5, 3},
	}

	ranks := RankFindNormalizedFold("limó", target)

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
	rs := Ranks{{"a", "b", 1, 0}, {"a", "cc", 2, 1}, {"a", "a", 0, 2}}
	wanted := Ranks{rs[2], rs[0], rs[1]}

	sort.Sort(rs)

	for i := range wanted {
		if wanted[i] != rs[i] {
			t.Errorf("expected %+v, got %+v", wanted, rs)
		}
	}
}

func BenchmarkMatch(b *testing.B) {
	ft := fuzzyTests[2]
	for b.Loop() {
		Match(ft.source, ft.target)
	}
}

func BenchmarkMatchBigLate(b *testing.B) {
	ft := fuzzyTests[0]
	for b.Loop() {
		Match(ft.source, ft.target)
	}
}

func BenchmarkMatchBigEarly(b *testing.B) {
	ft := fuzzyTests[1]
	for b.Loop() {
		Match(ft.source, ft.target)
	}
}

func BenchmarkMatchFold(b *testing.B) {
	ft := fuzzyTests[2]
	for b.Loop() {
		MatchFold(ft.source, ft.target)
	}
}

func BenchmarkMatchFoldBigLate(b *testing.B) {
	ft := fuzzyTests[0]
	for b.Loop() {
		MatchFold(ft.source, ft.target)
	}
}

func BenchmarkMatchFoldBigEarly(b *testing.B) {
	ft := fuzzyTests[1]
	for b.Loop() {
		MatchFold(ft.source, ft.target)
	}
}

func BenchmarkFindFold(b *testing.B) {
	b.Run("Plain", func(b *testing.B) { benchmarkFindFold(b, fuzzyTests[2]) })
	b.Run("BigLate", func(b *testing.B) { benchmarkFindFold(b, fuzzyTests[0]) })
	b.Run("BigEarly", func(b *testing.B) { benchmarkFindFold(b, fuzzyTests[1]) })
}

func benchmarkFindFold(b *testing.B, ft fuzzyTest) {
	src := ft.source
	var tgts []string
	for range 128 {
		tgts = append(tgts, ft.target)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FindFold(src, tgts)
	}
}

func BenchmarkRankMatch(b *testing.B) {
	ft := fuzzyTests[2]
	for b.Loop() {
		RankMatch(ft.source, ft.target)
	}
}

func BenchmarkRankMatchBigLate(b *testing.B) {
	ft := fuzzyTests[0]
	for b.Loop() {
		RankMatch(ft.source, ft.target)
	}
}

func BenchmarkRankMatchBigEarly(b *testing.B) {
	ft := fuzzyTests[1]
	for b.Loop() {
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
	// Output: [{Source:whl Target:cartwheel Distance:6 OriginalIndex:0} {Source:whl Target:wheel Distance:2 OriginalIndex:2}]
}

func TestRankMatchTransformedLength(t *testing.T) {
	for _, test := range []struct {
		name           string
		rank           func(string, string) int
		source, target string
		want           int
	}{
		{"normalized equal", RankMatchNormalized, "é", "e", 0},
		{"normalized subsequence", RankMatchNormalized, "éé", "eze", 1},
		{"normalized combining", RankMatchNormalized, "e\u0301", "e", 0},
		{"folded equal", RankMatchFold, "K", "k", 0},
		{"folded subsequence", RankMatchFold, "KK", "kzk", 1},
		{"normalized folded", RankMatchNormalizedFold, "ÉÉ", "eze", 1},
		{"not a match", RankMatchNormalized, "éé", "e", -1},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := test.rank(test.source, test.target); got != test.want {
				t.Errorf("rank(%q, %q) = %d, want %d", test.source, test.target, got, test.want)
			}
		})
	}
}

func ExampleMatchIndices() {
	fmt.Print(MatchIndices("twl", "cartwheel"))
	// Output: [3 4 8]
}

func BenchmarkMatchIndices(b *testing.B) {
	ft := fuzzyTests[2]
	for b.Loop() {
		MatchIndices(ft.source, ft.target)
	}
}

func BenchmarkMatchIndicesFold(b *testing.B) {
	ft := fuzzyTests[2]
	for b.Loop() {
		MatchIndicesFold(ft.source, ft.target)
	}
}
