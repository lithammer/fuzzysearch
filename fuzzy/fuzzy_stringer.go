// Fuzzy searching allows for flexibly matching a string with partial input,
// useful for filtering data very quickly based on lightweight user input.
package fuzzy

import (
	"fmt"
	"unicode/utf8"

	"golang.org/x/text/transform"
)

// MatchStringer returns true if source matches target using a fuzzy-searching
// algorithm. Note that it doesn't implement Levenshtein distance (see
// RankMatchStringer instead), but rather a simplified version where there's no
// approximation. The method will return true only if each character in the
// source can be found in the target and occurs after the preceding matches.
// This is the Stringer version of Match.
func MatchStringer(source string, target fmt.Stringer) bool {
	return matchStringer(source, target, noopTransformer())
}

// MatchFoldStringer is a case-insensitive version of MatchStringer.
func MatchFoldStringer(source string, target fmt.Stringer) bool {
	return matchStringer(source, target, foldTransformer())
}

// MatchNormalizedStringer is a unicode-normalized version of MatchStringer.
func MatchNormalizedStringer(source string, target fmt.Stringer) bool {
	return matchStringer(source, target, normalizeTransformer())
}

// MatchNormalizedFoldStringer is a unicode-normalized and case-insensitive version of MatchStringer.
func MatchNormalizedFoldStringer(source string, target fmt.Stringer) bool {
	return matchStringer(source, target, normalizedFoldTransformer())
}

func matchStringer(source string, target1 fmt.Stringer, transformer transform.Transformer) bool {
	source = stringTransform(source, transformer)
	target := stringTransform(target1.String(), transformer)

	lenDiff := len(target) - len(source)

	if lenDiff < 0 {
		return false
	}

	if lenDiff == 0 && source == target {
		return true
	}

Outer:
	for _, r1 := range source {
		for i, r2 := range target {
			if r1 == r2 {
				target = target[i+utf8.RuneLen(r2):]
				continue Outer
			}
		}
		return false
	}

	return true
}

// FindStringer will return a list of strings in targets that fuzzy matches source.
func FindStringer(source string, targets []fmt.Stringer) []fmt.Stringer {
	return findStringer(source, targets, noopTransformer())
}

// FindFoldStringer is a case-insensitive version of FindStringer.
func FindFoldStringer(source string, targets []fmt.Stringer) []fmt.Stringer {
	return findStringer(source, targets, foldTransformer())
}

// FindNormalizedStringer is a unicode-normalized version of FindStringer.
func FindNormalizedStringer(source string, targets []fmt.Stringer) []fmt.Stringer {
	return findStringer(source, targets, normalizeTransformer())
}

// FindNormalizedFoldStringer is a unicode-normalized and case-insensitive version of FindStringer.
func FindNormalizedFoldStringer(source string, targets []fmt.Stringer) []fmt.Stringer {
	return findStringer(source, targets, normalizedFoldTransformer())
}

func findStringer(source string, targets []fmt.Stringer, transformer transform.Transformer) []fmt.Stringer {
	var matches []fmt.Stringer

	for _, target := range targets {
		if matchStringer(source, target, transformer) {
			matches = append(matches, target)
		}
	}

	return matches
}

// RankMatchStringer is similar to MatchStringer except it will measure the Levenshtein
// distance between the source and the target and return its result. If there
// was no matchStringer, it will return -1.
// Given the requirements of matchStringer, RankMatchStringer only needs to perform a subset of
// the Levenshtein calculation, only deletions need be considered, required
// additions and substitutions would fail the matchStringer test.
// This is the Stringer version of RankMatch.
func RankMatchStringer(source string, target fmt.Stringer) int {
	return rankStringer(source, target, noopTransformer())
}

// RankMatchFoldStringer is a case-insensitive version of RankMatchStringer.
func RankMatchFoldStringer(source string, target fmt.Stringer) int {
	return rankStringer(source, target, foldTransformer())
}

// RankMatchNormalizedStringer is a unicode-normalized version of RankMatchStringer.
func RankMatchNormalizedStringer(source string, target fmt.Stringer) int {
	return rankStringer(source, target, normalizeTransformer())
}

// RankMatchNormalizedFoldStringer is a unicode-normalized and case-insensitive version of RankMatchStringer.
func RankMatchNormalizedFoldStringer(source string, target fmt.Stringer) int {
	return rankStringer(source, target, normalizedFoldTransformer())
}

func rankStringer(source string, target1 fmt.Stringer, transformer transform.Transformer) int {
	lenDiff := len(target1.String()) - len(source)

	if lenDiff < 0 {
		return -1
	}

	source = stringTransform(source, transformer)
	target := stringTransform(target1.String(), transformer)

	if lenDiff == 0 && source == target {
		return 0
	}

	runeDiff := 0

Outer:
	for _, r1 := range source {
		for i, r2 := range target {
			if r1 == r2 {
				target = target[i+utf8.RuneLen(r2):]
				continue Outer
			} else {
				runeDiff++
			}
		}
		return -1
	}

	// Count up remaining char
	runeDiff += utf8.RuneCountInString(target)

	return runeDiff
}

// RankFindStringer is similar to FindStringer, except it will also rankStringer all matches using
// Levenshtein distance.
func RankFindStringer(source string, targets []fmt.Stringer) RanksStringer {
	return rankFindStringer(source, targets, noopTransformer())
}

// RankFindFoldStringer is a case-insensitive version of RankFindStringer.
func RankFindFoldStringer(source string, targets []fmt.Stringer) RanksStringer {
	return rankFindStringer(source, targets, foldTransformer())
}

// RankFindNormalizedStringer is a unicode-normalized version of RankFindStringer.
func RankFindNormalizedStringer(source string, targets []fmt.Stringer) RanksStringer {
	return rankFindStringer(source, targets, normalizeTransformer())
}

// RankFindNormalizedFoldStringer is a unicode-normalized and case-insensitive version of RankFindStringer.
func RankFindNormalizedFoldStringer(source string, targets []fmt.Stringer) RanksStringer {
	return rankFindStringer(source, targets, normalizedFoldTransformer())
}

func rankFindStringer(source string, targets []fmt.Stringer, transformer transform.Transformer) RanksStringer {
	var r RanksStringer

	for index, target := range targets {
		if matchStringer(source, target, transformer) {
			distance := LevenshteinDistance(source, target.String())
			r = append(r, RankStringer{source, target, distance, index})
		}
	}
	return r
}

type RankStringer struct {
	// Source is used as the source for matching.
	Source string

	// Target is the word matched against.
	Target fmt.Stringer

	// Distance is the Levenshtein distance between Source and Target.
	Distance int

	// Location of Target in original list
	OriginalIndex int
}

type RanksStringer []RankStringer

func (r RanksStringer) Len() int {
	return len(r)
}

func (r RanksStringer) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r RanksStringer) Less(i, j int) bool {
	return r[i].Distance < r[j].Distance
}
