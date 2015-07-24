# Fuzzy Search

[![Build Status](https://travis-ci.org/renstrom/fuzzysearch.svg?branch=master)](https://travis-ci.org/renstrom/fuzzysearch)
[![Godoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/renstrom/fuzzysearch/fuzzy)

A Go port of _[bevacqua/fuzzysearch][1]_ (written in JavaScript).

Fuzzy searching allows for flexibly matching a string with partial input, useful for filtering data very quickly based on lightweight user input.

Returns `true` if `needle` matches `haystack` using a fuzzy-searching algorithm. Note that this program doesn't implement _[Levenshtein distance][2]_, but rather a simplified version where **there's no approximation**. The method will return `true` only if each character in the `needle` can be found in the `haystack` and occurs after the preceding matches.

The current implementation uses the algorithm suggested by Mr. Aleph, a russian compiler engineer working at V8.

## Differences to bevacqua/fuzzysearch

Includes `SearchMany()`, which is a convenience function to help filter a list of words which mathes.

Also, in contrast to the original JavaScript implementation, this library supplies two extra functions to help with ranking matches using Levenshtein distance. Namely `RankSearch()` and `RankSearchMany()`.

## Usage

```go
fuzzy.Search("twl", "cartwheel")  // true
fuzzy.Search("cart", "cartwheel") // true
fuzzy.Search("cw", "cartwheel")   // true
fuzzy.Search("ee", "cartwheel")   // true
fuzzy.Search("art", "cartwheel")  // true
fuzzy.Search("eeel", "cartwheel") // false
fuzzy.Search("dog", "cartwheel")  // false

fuzzy.RankSearch("kitten", "sitten") // 3

words := []string{"cartwheel", "foobar", "wheel", "baz"})
fuzzy.SearchMany("whl", words) // [cartwheel, wheel]

fuzzy.RankSearchMany("whl", words) // [{cartwheel 6} {wheel 2}]
```

## License

MIT

[1]: https://github.com/bevacqua/fuzzysearch
[2]: http://en.wikipedia.org/wiki/Levenshtein_distance
