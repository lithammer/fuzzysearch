# Fuzzy Search

[![Build Status](https://travis-ci.org/renstrom/fuzzysearch.svg?branch=master)](https://travis-ci.org/renstrom/fuzzysearch)
[![Godoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/renstrom/fuzzysearch)

A Go port of _[bevacqua's/fuzzysearch][1]_ written in JavaScript.

Fuzzy searching allows for flexibly matching a string with partial input, useful for filtering data very quickly based on lightweight user input.

Returns `true` if `needle` matches `haystack` using a fuzzy-searching algorithm. Note that this program doesn't implement _[levenshtein distance][2]_, but rather a simplified version where **there's no approximation**. The method will return `true` only if each character in the `needle` can be found in the `haystack` and occurs after the preceding matches.

The current implementation uses the algorithm suggested by Mr. Aleph, a russian compiler engineer working at V8.

## Usage

```go
import "github.com/renstrom/fuzzysearch/fuzzy"

func main() {
    fuzzy.Search("twl", "cartwheel")  // <- true
    fuzzy.Search("cart", "cartwheel") // <- true
    fuzzy.Search("cw", "cartwheel")   // <- true
    fuzzy.Search("ee", "cartwheel")   // <- true
    fuzzy.Search("art", "cartwheel")  // <- true
    fuzzy.Search("eeel", "cartwheel") // <- false
    fuzzy.Search("dog", "cartwheel")  // <- false
}
```

## License

MIT

[1]: https://github.com/bevacqua/fuzzysearch
[2]: http://en.wikipedia.org/wiki/Levenshtein_distance
