package stringhelper

import (
	"strings"
	"unicode/utf8"
)

// ******** Public functions ********

// SplitAny splits a string at any character in the separators string.
func SplitAny(source string, separators string) []string {
	if len(separators) == 0 {
		return explode(source)
	}

	result := make([]string, countAny(source, separators)+1)

	i := 0
	for {
		pos := strings.IndexAny(source, separators)
		if pos == -1 {
			break
		}

		result[i] = source[:pos]
		source = source[pos+1:]
		i++
	}

	result[i] = source

	return result
}

// CountAny counts all occurrences of any character in the separators string.
func CountAny(source string, separators string) int {
	if len(separators) == 0 || len(source) == 0 {
		return -1
	}

	return countAny(source, separators)
}

// ******** Private functions ********

// countAny counts all occurrences of any character in the separators string.
// This is the internal version without any checks.
func countAny(source string, separators string) int {
	n := 0
	for {
		pos := strings.IndexAny(source, separators)
		if pos == -1 {
			break
		}

		n++
		source = source[pos+1:]
	}

	return n
}

// explode splits s into a slice of UTF-8 strings, one string per Unicode character.
// Invalid UTF-8 bytes are sliced individually.
func explode(source string) []string {
	runeCount := utf8.RuneCountInString(source)
	result := make([]string, runeCount)
	for i := 0; i < runeCount; i++ {
		_, size := utf8.DecodeRuneInString(source)
		result[i] = source[:size]
		source = source[size:]
	}

	return result
}
