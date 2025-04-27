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
// Version: 5.0.0
//
// Change history:
//    2025-03-12: V1.0.0: Created.
//    2025-03-15: V2.0.0: Parallelize decryption.
//    2025-03-17: V2.1.0: Use clear.
//    2025-03-23: V3.0.0: Refactored interface.
//    2025-03-23: V4.0.0: Make generic.
//    2025-04-27: V5.0.0: Use object.
//

package transposition

import (
	"sync"
)

// ******** Public functions ********

// Untranspose untransposes a slice.
// Side effects: Source is overwritten if there is more than one password.
// The list of passwords is reversed.
func (r *Transposition[T]) Untranspose(source []T) []T {
	result := make([]T, len(source))

	from := result
	to := source
	maxIndex := len(r.orders) - 1
	for i := range r.orders {
		from, to = to, from

		untransposeToTarget(to, from, r.orders[maxIndex-i])
	}

	return to
}

// ******** Private functions ********

// untransposeToTarget reverts a transposition with the given password to the given target.
func untransposeToTarget[T any](target []T, source []T, offsets []int) {
	sourceLen := len(source)

	transposeLen := len(offsets)

	var wg sync.WaitGroup
	sourceIndex := 0
	for _, offset := range offsets {
		// Untranspose each column in parallel.
		wg.Add(1)
		go untransposeColumn(&wg, target, source, sourceLen, transposeLen, offset, sourceIndex)

		sourceIndex += columnLen(sourceLen, transposeLen, offset)
	}

	wg.Wait()
}

// untransposeColumn untransposes a column.
func untransposeColumn[T any](
	wg *sync.WaitGroup,
	to []T,
	from []T,
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
