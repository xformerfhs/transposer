//
// SPDX-FileCopyrightText: Copyright 2025 Frank Schwab
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileType: SOURCE
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
//
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Author: Frank Schwab
//
// Version: 1.0.0
//
// Change history:
//    2025-03-23: V1.0.0: Created.
//

package stringhelper

import (
	"iter"
	"strings"
	"unicode/utf8"
)

// ******** Public functions ********

// SplitAny splits a string at any character in the separators string.
func SplitAny(source string, separators string) []string {
	return genSplitAny(source, separators, -1)
}

// SplitAnyN splits a string at any character in the separators string.
// It returns at most n substrings; the last substring will be the unsplit remainder.
func SplitAnyN(source string, separators string, n int) []string {
	return genSplitAny(source, separators, n)
}

// SplitAnySeq returns an iterator over all substrings of s separated by separator.
// The iterator yields the same strings that would be returned by [SplitAny](source, separators),
// but without constructing the slice.
// It returns a single-use iterator.
func SplitAnySeq(source string, separators string) iter.Seq[string] {
	if len(separators) == 0 {
		return explodeSeq(source)
	}

	return func(yield func(string) bool) {
		for {
			pos := strings.IndexAny(source, separators)
			if pos < 0 {
				break
			}
			if !yield(source[:pos]) {
				return
			}

			// size should have been returned by IndexAny.
			_, size := utf8.DecodeRuneInString(source[pos:])
			source = source[pos+size:]
		}

		yield(source)
	}
}

// CountAny counts all occurrences of any character in the separators string.
func CountAny(source string, separators string) int {
	if len(separators) == 0 {
		return utf8.RuneCountInString(source) + 1
	}

	return countAny(source, separators)
}

// ******** Private functions ********

// genSplitAny splits after each instance of sep.
// It returns at most n substrings; the last substring will be the unsplit remainder.
func genSplitAny(source string, separators string, n int) []string {
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

		// size should have been returned by IndexAny.
		_, size := utf8.DecodeRuneInString(source[pos:])

		result[i] = source[:pos]
		source = source[pos+size:]

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

		// size should have been returned by IndexAny.
		_, size := utf8.DecodeRuneInString(source[pos:])
		source = source[pos+size:]
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

// explodeSeq returns an iterator over the runes in s.
func explodeSeq(source string) iter.Seq[string] {
	return func(yield func(string) bool) {
		for len(source) > 0 {
			_, size := utf8.DecodeRuneInString(source)
			if !yield(source[:size]) {
				return
			}
			source = source[size:]
		}
	}
}
