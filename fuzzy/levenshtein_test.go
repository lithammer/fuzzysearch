package fuzzy

import "testing"

var levenshteinDistanceTests = []struct {
	s, t   string
	wanted int
}{
	{"stakes", "The only really major change from the <em>status quo ante</em> resulting from the War of Austrian Succession was Austria’s loss of Silesia to Prussia. As a consequence, Frederick increased the number of his subjects by over 50 percent and made Prussia the first new entrant to the ranks of the Great Powers since the Peace of Westphalia in 1648. Frederick—who by this point was already being called “the Great” by his admirers—had demonstrated that as a political instrument war could accomplish major objectives even in an age of “limited war,” when sociopolitical constraints as well as military ones tended to make wars relatively indecisive and profitless. Still, the Prussian king was able to make such substantial gains only because his superior soldiers and his increasingly adept generalship enabled him to win five consecutive battles against the Austrians. Even that impressive military record, however, bought Prussia only an insecure conquest and a cold and temporary peace. Maria Theresa remained determined to recover Silesia, and worked tirelessly to improve her army and adjust the political situation to make that possible. The result, the subject of chapter 9, was a renewed effort to undo Frederick’s gains, a war that would push Austria, France, Britain, and Prussia towards a style of warfare employing greater resources and involving higher political stakes.", 1374},
	{"a", "a", 0},
	{"ab", "ab", 0},
	{"ab", "aa", 1},
	{"ab", "aa", 1},
	{"ab", "aaa", 2},
	{"bbb", "a", 3},
	{"kitten", "sitting", 3},
	{"ёлка", "ёлочка", 2},
	{"ветер", "ёлочка", 6},
	{"中国", "中华人民共和国", 5},
	{"日本", "中华人民共和国", 7},
}

func TestLevenshtein(t *testing.T) {
	for _, test := range levenshteinDistanceTests {
		distance := LevenshteinDistance(test.s, test.t)
		if distance != test.wanted {
			t.Errorf("got distance %d, expected %d for %s in %s", distance, test.wanted, test.s, test.t)
		}
	}
}

func BenchmarkLevenshteinDistance(b *testing.B) {
	ldt := levenshteinDistanceTests[1]
	ldt2 := levenshteinDistanceTests[5]
	for i := 0; i < b.N; i++ {
		LevenshteinDistance(ldt.s, ldt.t)
		LevenshteinDistance(ldt2.s, ldt2.t)
	}
}

func BenchmarkLevenshteinDistanceBig(b *testing.B) {
	ldt := levenshteinDistanceTests[0]
	for i := 0; i < b.N; i++ {
		LevenshteinDistance(ldt.s, ldt.t)
	}
}
