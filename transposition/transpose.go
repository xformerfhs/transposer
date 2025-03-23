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
// Version: 4.0.0
//
// Change history:
//    2025-03-12: V1.0.0: Created.
//    2025-03-15: V2.0.0: Parallelize encryption.
//    2025-03-17: V2.1.0: Use clear.
//    2025-03-23: V3.0.0: Refactored interface.
//    2025-03-23: V4.0.0: Make generic.
//

package transposition

import (
	"sync"
)

// ******** Public functions ********

// TransposeToTarget transposes a slice with the given password to the given target.
func TransposeToTarget[S ~[]T, T any](target S, source S, password string) {
	sourceLen := len(source)

	offsets := columnOrder(password)
	transposeLen := len(offsets)

	var wg sync.WaitGroup
	destinationIndex := 0
	for _, offset := range offsets {
		// Transpose each column in parallel.
		wg.Add(1)
		go transposeColumn(&wg, target, source, sourceLen, transposeLen, offset, destinationIndex)

		destinationIndex += columnLen(sourceLen, transposeLen, offset)
	}

	wg.Wait()

	clear(offsets)
}

// Transpose transposes a slice with the given password.
func Transpose[S ~[]T, T any](source S, password string) S {
	result := make(S, len(source))
	TransposeToTarget(result, source, password)
	return result
}

// TransposeMultiplePasswords transposes a rune array with multiple passwords.
// Side effect: Source is overwritten, if there is more than one password.
func TransposeMultiplePasswords[S ~[]T, T any](source S, passwords []string) S {
	result := make(S, len(source))
	from := result
	to := source
	for _, password := range passwords {
		from, to = to, from

		TransposeToTarget(to, from, password)
	}

	return to
}

// ******** Private functions ********

// transposeColumn transposes a column.
func transposeColumn[S ~[]T, T any](
	wg *sync.WaitGroup,
	to S,
	from S,
	sourceLen int,
	transposeLen int,
	offset int,
	destinationIndex int) {
	defer wg.Done()

	for sourceIndex := offset; sourceIndex < sourceLen; sourceIndex += transposeLen {
		to[destinationIndex] = from[sourceIndex]
		destinationIndex++
	}
}
