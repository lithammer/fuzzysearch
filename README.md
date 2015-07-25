# Fuzzy Search

[![Build Status](https://travis-ci.org/renstrom/fuzzysearch.svg?branch=master)](https://travis-ci.org/renstrom/fuzzysearch)
[![Godoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/renstrom/fuzzysearch/fuzzy)

A Go port of _[bevacqua/fuzzysearch][1]_ (written in JavaScript).

Fuzzy searching allows for flexibly matching a string with partial input, useful for filtering data very quickly based on lightweight user input.

Returns `true` if `needle` matches `haystack` using a fuzzy-searching algorithm. Note that this program doesn't implement _[Levenshtein distance][2]_, but rather a simplified version where **there's no approximation**. The method will return `true` only if each character in the `needle` can be found in the `haystack` and occurs after the preceding matches.

The current implementation uses the algorithm suggested by Mr. Aleph, a russian compiler engineer working at V8.

## Differences to bevacqua/fuzzysearch

The function `search()` has been renamed to `Match()`.

Includes `Find()`, which is a convenience function to help filter a list of words which mathes.

Also, in contrast to the original JavaScript implementation, this library supplies two extra functions to help with ranking matches using Levenshtein distance. Namely `RankMatch()` and `RankFind()`.

## Usage

```go
fuzzy.Match("twl", "cartwheel")  // true
fuzzy.Match("cart", "cartwheel") // true
fuzzy.Match("cw", "cartwheel")   // true
fuzzy.Match("ee", "cartwheel")   // true
fuzzy.Match("art", "cartwheel")  // true
fuzzy.Match("eeel", "cartwheel") // false
fuzzy.Match("dog", "cartwheel")  // false

fuzzy.RankMatch("kitten", "sitting") // 3

words := []string{"cartwheel", "foobar", "wheel", "baz"})
fuzzy.Find("whl", words) // [cartwheel wheel]

fuzzy.RankFind("whl", words) // [{cartwheel 6} {wheel 2}]
```

## License

MIT

[1]: https://github.com/bevacqua/fuzzysearch
[2]: http://en.wikipedia.org/wiki/Levenshtein_distance
