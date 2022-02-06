# Fuzzy Search

Inspired by [bevacqua/fuzzysearch][1], a fuzzy matching library written in
JavaScript. But contains some extras like ranking using [Levenshtein
distance][2] and finding matches in a list of words.

Fuzzy searching allows for flexibly matching a string with partial input,
useful for filtering data very quickly based on lightweight user input.

The current implementation uses the algorithm suggested by Mr. Aleph, a russian
compiler engineer working at V8.

## Install

```
go get github.com/lithammer/fuzzysearch/fuzzy
```

## Usage

```go
package main

import "github.com/lithammer/fuzzysearch/fuzzy"

func main() {
	fuzzy.Match("twl", "cartwheel")  // true
	fuzzy.Match("cart", "cartwheel") // true
	fuzzy.Match("cw", "cartwheel")   // true
	fuzzy.Match("ee", "cartwheel")   // true
	fuzzy.Match("art", "cartwheel")  // true
	fuzzy.Match("eeel", "cartwheel") // false
	fuzzy.Match("dog", "cartwheel")  // false
	fuzzy.Match("kitten", "sitting") // false
	
	fuzzy.RankMatch("kitten", "sitting") // -1
	fuzzy.RankMatch("cart", "cartwheel") // 5
	
	words := []string{"cartwheel", "foobar", "wheel", "baz"}
	fuzzy.Find("whl", words) // [cartwheel wheel]
	
	fuzzy.RankFind("whl", words) // [{whl cartwheel 6 0} {whl wheel 2 2}]
	
	// Unicode normalized matching.
	fuzzy.MatchNormalized("cartwheel", "cartwhéél") // true

	// Case insensitive matching.
	fuzzy.MatchFold("ArTeeL", "cartwheel") // true
}
```

You can sort the result of a `fuzzy.RankFind()` call using the [`sort`][3]
package in the standard library:

```go
matches := fuzzy.RankFind("whl", words) // [{whl cartwheel 6 0} {whl wheel 2 2}]
sort.Sort(matches) // [{whl wheel 2 2} {whl cartwheel 6 0}]
```

See the [`fuzzy`][4] package documentation for more examples.

## License

MIT

[1]: https://github.com/bevacqua/fuzzysearch
[2]: http://en.wikipedia.org/wiki/Levenshtein_distance
[3]: https://golang.org/pkg/sort/
[4]: https://pkg.go.dev/github.com/lithammer/fuzzysearch/fuzzy
