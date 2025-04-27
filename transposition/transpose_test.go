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
// Version: 1.1.0
//
// Change history:
//    2025-03-16: V1.0.0: Created.
//    2025-03-23: V1.1.0: Adapt to new interface.
//

// Package transposition_test contains tests for transpositions.
package transposition_test

import (
	"math/rand/v2"
	"slices"
	"testing"
	"transposer/transposition"
)

// ******** Private constants ********

// fmtExpectedActual is the format expected and actual rune slice values.
const fmtExpectedActual = "\nExpected: %c\nActual: %c\n"

// pwLengths is the list of possible password lengths.
var pwLengths = []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31}

// sourceCharacters is the list of characters a source is constructed from
var sourceCharacters = []rune(`ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÜabcdefghijklmnopqrstuvwxyzäöüß0123456789É`)

// expected100 is the expected transposition of a source of 100 characters.
var expected100 = []rune(`rMFlHAÖUwLEpF9KDkE8fÖzPIjIBdZTMtNGoA4YRsRKmB5bXßZSqGÉaW3XQiD7cY2SLnC6gÜ0ÜVüWPxJCeÄöÄTväVOuQJha1UNyOH`)

// knownPasswords is the list of passwords for the test with 100 source characters.
var knownPasswords = []string{`galadriel`, `legolas`}

// ******** Helper functions *********

// buildSource builds a source with the given length.
func buildSource(l int) []rune {
	result := make([]rune, l)

	j := 0
	for i := 0; i < l; i++ {
		result[i] = sourceCharacters[j]

		j = (j + 1) % len(sourceCharacters)
	}

	return result
}

// buildRandomPassword builds a random password of the given length.
func buildRandomPassword(l int) string {
	result := make([]rune, l)
	for i := 0; i < l; i++ {
		result[i] = 'A' + rand.Int32N(26)
	}

	return string(result)
}

// ******** Test functions ********

// TestTransposeKnown transposes a known source and compares the result with the known result.
func TestTransposeKnown(t *testing.T) {
	source := buildSource(100)
	transpose := transposition.New[rune](knownPasswords)
	transposed := transpose.Transpose(source)
	if !slices.Equal(transposed, expected100) {
		t.Fatalf(fmtExpectedActual, expected100, transposed)
	}
	transpose.Destroy()
}

// TestUntransposeKnown transposes a known transposed source and compares the result with the known source.
func TestUntransposeKnown(t *testing.T) {
	expected := buildSource(100)
	transpose := transposition.New[rune](knownPasswords)
	transposed := transpose.Untranspose(expected100)
	if !slices.Equal(transposed, expected) {
		t.Fatalf(fmtExpectedActual, expected100, transposed)
	}
	transpose.Destroy()
}

// TestTransposeRandom transposes a random source with a random number of random length passwords
// and compares the untransposed value of those with the original.
func TestTransposeRandom(t *testing.T) {
	for i := 0; i < 100; i++ {
		numPasswords := rand.IntN(10) + 1
		passwords := make([]string, numPasswords)
		for j := 0; j < numPasswords; j++ {
			passwords[j] = buildRandomPassword(pwLengths[rand.IntN(len(pwLengths))])
		}

		source := buildSource(rand.IntN(180) + 20)
		// Transpose overwrites the source, so it must be saved here.
		sourceSafe := make([]rune, len(source))
		copy(sourceSafe, source)

		transpose := transposition.New[rune](passwords)
		transposed := transpose.Transpose(source)
		untransposed := transpose.Untranspose(transposed)

		if !slices.Equal(untransposed, sourceSafe) {
			t.Fatalf(fmtExpectedActual, sourceSafe, untransposed)
		}

		transpose.Destroy()
	}
}
