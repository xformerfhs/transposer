package stringhelper

import (
	"strings"
	"unicode/utf8"
)

// ******** Public functions ********

// SplitAny splits a string at any character in the separators string.
func SplitAny(source string, separators string) []string {
	return genSplitAny(source, separators, 0, -1)
}

// SplitAnyN splits a string at any character in the separators string.
// It returns at most n substrings; the last substring will be the unsplit remainder.
func SplitAnyN(source string, separators string, n int) []string {
	return genSplitAny(source, separators, 0, n)
}

// CountAny counts all occurrences of any character in the separators string.
func CountAny(source string, separators string) int {
	if len(separators) == 0 {
		return utf8.RuneCountInString(source) + 1
	}

	return countAny(source, separators)
}

// ******** Private functions ********

// Generic split: splits after each instance of sep,
// including sepSave bytes of sep in the subarrays.
func genSplitAny(source string, separators string, sepSave int, n int) []string {
	if n == 0 {
		return nil
	}

	if len(separators) == 0 {
		return explodeN(source, n)
	}

	if n < 0 {
		n = countAny(source, separators) + 1
	}

	if n > len(source)+1 {
		n = len(source) + 1
	}

	result := make([]string, n)
	n--
	i := 0
	for i < n {
		pos := strings.IndexAny(source, separators)
		if pos < 0 {
			break
		}

		result[i] = source[:pos+sepSave]
		source = source[pos+1:]

		i++
	}

	result[i] = source

	return result[:i+1]
}

// countAny counts all occurrences of any character in the separators string.
// This is the internal version without any checks.
func countAny(source string, separators string) int {
	n := 0
	for {
		pos := strings.IndexAny(source, separators)
		if pos < 0 {
			break
		}

		n++
		source = source[pos+1:]
	}

	return n
}

// explodeN splits source into a slice of UTF-8 strings,
// one string per Unicode character up to a maximum of n (n < 0 means no limit).
// Invalid UTF-8 bytes are sliced individually.
func explodeN(source string, n int) []string {
	runeCount := utf8.RuneCountInString(source)
	if n < 0 || n > runeCount {
		n = runeCount
	}

	result := make([]string, n)
	for i := 0; i < n-1; i++ {
		_, size := utf8.DecodeRuneInString(source)
		result[i] = source[:size]
		source = source[size:]
	}

	if n > 0 {
		result[n-1] = source
	}

	return result
}
