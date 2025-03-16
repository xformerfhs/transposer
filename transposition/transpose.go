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
// Version: 2.0.0
//
// Change history:
//    2025-03-12: V1.0.0: Created.
//    2025-03-15: V2.0.0: Parallelize encryption.
//

package transposition

import (
	"sync"
	"transposer/slicehelper"
)

// ******** Public functions ********

// TransposeRuneArray transposes a rune array with the given passwords.
// Attention: source is overwritten, if there is more than one password!
func TransposeRuneArray(source []rune, passwords []string) []rune {
	sourceLen := len(source)
	result := make([]rune, sourceLen)

	from := result
	to := source
	for _, password := range passwords {
		from, to = to, from

		offsets := columnOrder(password)
		transposeLen := len(offsets)

		var wg sync.WaitGroup
		destinationIndex := 0
		for _, offset := range offsets {
			// Transpose each column in parallel.
			wg.Add(1)
			go transposeColumn(&wg, from, to, sourceLen, transposeLen, offset, destinationIndex)

			destinationIndex += columnLen(sourceLen, transposeLen, offset)
		}

		wg.Wait()

		slicehelper.ClearInteger(offsets)
	}

	return to
}

// ******** Private functions ********

// transposeColumn transposes a column.
func transposeColumn(
	wg *sync.WaitGroup,
	from []rune,
	to []rune,
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
