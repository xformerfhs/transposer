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
//    2025-04-27: V1.0.0: Created.
//

package transposition

// ******** Public types ********

// Transposition transposes slices.
type Transposition[T any] struct {
	orders [][]int
}

// ******** Public creation and destruction functions ********

// New creates a new transposition object.
func New[T any](passwords []string) *Transposition[T] {
	orders := make([][]int, len(passwords))
	for i, password := range passwords {
		orders[i] = columnOrder(password)
	}

	return &Transposition[T]{
		orders: orders,
	}
}

// Destroy destroys the transposition object and removes the secret order information from memory.
// After calling Destroy, the transposition object is no longer usable.
func (r *Transposition[T]) Destroy() {
	for i := range r.orders {
		clear(r.orders[i])
		r.orders[i] = nil
	}

	r.orders = nil
}
