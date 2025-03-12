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
//    2023-05-03: V1.0.0: Created. fs
//    2025-02-22: V2.0.0: More efficient Fill. Add FillToCap. fs

// Package slicehelper contains useful functions for slices that are missing from
// the standard packages.
package slicehelper

import (
	"transposer/constraints"
)

// ******** Private constants ********

// powerFillThresholdLen is the slice length where PowerFill is more efficient than SimpleFill.
const powerFillThresholdLen = 74

// ******** Public functions ********

// Fill fills a slice with a value in an efficient way up to its length.
func Fill[S ~[]T, T any](s S, v T) {
	sLen := len(s)

	if sLen > 0 {
		doFill(s, v, sLen)
	}
}

// FillToCap fills a slice with a value in an efficient way up to its capacity.
func FillToCap[S ~[]T, T any](s S, v T) {
	sLen := cap(s)

	if sLen > 0 {
		doFill(s[:sLen], v, sLen)
	}
}

// ClearInteger clears an integer type slice.
func ClearInteger[S ~[]T, T constraints.Integer](a S) {
	FillToCap(a, 0)
}

// SafeClearInteger clears an integer type slice, without worrying if it is nil
func SafeClearInteger[S ~[]T, T constraints.Integer](a S) {
	if a != nil {
		FillToCap(a, 0)
	}
}

// CopyOfWithLen returns a copy of slice s with length l.
// l may be longer or shorter than the length of s.
func CopyOfWithLen[S ~[]T, T any](s S, l int) []T {
	r := make(S, l)
	copy(r, s)

	return r
}

// Concat returns a new slice concatenating the passed in slices.
// This is a streamlined version of the slices.Concat function of Go V1.22.
func Concat[S ~[]T, T any](slices ...S) S {
	// 1. Calculate total size.
	size := 0
	for _, s := range slices {
		size += len(s)
	}

	// 2. Make new slice with the total size as the capacity and 0 length.
	result := make(S, 0, size)

	// 3. Append all source slices.
	for _, s := range slices {
		result = append(result, s...)
	}

	return result
}

// SplitTail cuts the last l elements from the supplied slice.
func SplitTail[S ~[]T, T any](s S, l int) (S, S) {
	p := len(s) - l
	return s[:p], s[p:]
}

// ******** Private functions ********

// doFill fills a slice in an optimal way.
func doFill[S ~[]T, T any](s S, v T, l int) {
	if l >= powerFillThresholdLen {
		doPowerFill(s, v, l)
	} else {
		doSimpleFill(s, v, l)
	}
}

// doSimpleFill fills a slice in a simple way.
func doSimpleFill[S ~[]T, T any](s S, v T, l int) {
	for i := 0; i < l; i++ {
		s[i] = v
	}
}

// doPowerFill fills a slice in an efficient way.
func doPowerFill[S ~[]T, T any](s S, v T, l int) {
	// Put the value into the first slice element
	s[0] = v
	// Incrementally duplicate the value into the rest of the slice
	for j := 1; j < l; j <<= 1 {
		copy(s[j:], s[:j])
	}
}
