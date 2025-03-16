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
//    2025-03-15: V2.0.0: Parallelize decryption.
//

package transposition

import (
	"slices"
	"sync"
	"transposer/slicehelper"
)

// UnTransposeRuneArray reverts a transposition with the given passwords.
// Attention: source is overwritten, if there is more than one password!
func UnTransposeRuneArray(source []rune, passwords []string) []rune {
	sourceLen := len(source)
	result := make([]rune, sourceLen)

	from := result
	to := source
	// For decryption passwords have to be used last to first.
	slices.Reverse(passwords)
	for _, password := range passwords {
		from, to = to, from

		offsets := columnOrder(password)
		transposeLen := len(offsets)

		var wg sync.WaitGroup
		sourceIndex := 0
		for _, offset := range offsets {
			// Untranspose each column in parallel.
			wg.Add(1)
			go untransposeColumn(&wg, from, to, sourceLen, transposeLen, offset, sourceIndex)

			sourceIndex += columnLen(sourceLen, transposeLen, offset)
		}

		wg.Wait()

		slicehelper.ClearInteger(offsets)
	}

	return to
}

// ******** Private functions ********

// untransposeColumn untransposes a column.
func untransposeColumn(
	wg *sync.WaitGroup,
	from []rune,
	to []rune,
	sourceLen int,
	transposeLen int,
	offset int,
	sourceIndex int) {
	defer wg.Done()

	for destinationIndex := offset; destinationIndex < sourceLen; destinationIndex += transposeLen {
		to[destinationIndex] = from[sourceIndex]
		sourceIndex++
	}
}
