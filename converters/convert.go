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
//    2025-03-14: V1.0.0: Created.
//

// Package converters contains functions that convert one rune into another.
package converters

import "unicode"

// ******** Public types *********

// RuneConverter is a function that converts a rune into another one.
type RuneConverter func(rune) rune

// ******** Public functions *********

// Same returns the rune unchanged.
func Same(r rune) rune {
	return r
}

// ToLower converts letters to lower case, others are not changed.
func ToLower(r rune) rune {
	if unicode.IsLetter(r) {
		return unicode.ToLower(r)
	}

	return r
}

// ToUpper converts letters to upper case, others are not changed.
func ToUpper(r rune) rune {
	if unicode.IsLetter(r) {
		return unicode.ToUpper(r)
	}

	return r
}
