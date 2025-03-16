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
//    2025-03-15: V2.0.0: Added columnLen.
//

// Package transposition contains the functions that transpose and untranspose
// a rune slice.
package transposition

import "transposer/linkedlist"

// ******** Private functions ********

// columnOrder returns a slice with the order of the offsets of the columns.
func columnOrder(source string) []int {
	orderList := linkedlist.New()

	// Range over a string returns the *byte offset* and the rune value
	// So we have to use our own rune index
	i := 0

	for _, char := range source {
		orderList.Insert(i, char)
		i++
	}

	return orderList.ValueOrderedIndices()
}

// columnLen calculates the length of the column starting at the given offset.
func columnLen(sourceLen int, transposeLen int, offset int) int {
	sourceLen = sourceLen - offset
	result := sourceLen / transposeLen
	remainder := sourceLen - (result * transposeLen)
	if remainder > 0 {
		result++
	}

	return result
}
