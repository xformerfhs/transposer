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
//    2025-03-12: V1.0.0: Created.
//

package transposition

import (
	"transposer/slicehelper"
)

// ******** Public functions ********

// TransposeRuneArray transposes a rune array with the given passwords.
func TransposeRuneArray(source []rune, passwords []string) []rune {
	sourceLength := len(source)
	result := make([]rune, sourceLength)

	from := result
	to := source
	for _, password := range passwords {
		from, to = to, from

		offsets := columnOrder(password)
		transposeLength := len(offsets)

		destinationIndex := 0
		for _, offset := range offsets {
			for sourceIndex := offset; sourceIndex < sourceLength; sourceIndex += transposeLength {
				to[destinationIndex] = from[sourceIndex]
				destinationIndex++
			}
		}

		slicehelper.ClearInteger(offsets)
	}

	return to
}
