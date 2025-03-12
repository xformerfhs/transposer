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

// Package linkedlist contains functions to implement a doubly-linked list that
// contains nodes with a rune and an index.
package linkedlist

// ******** Public types ********

// SortedRuneList implements a sorted list of values.
type SortedRuneList struct {
	head   *node
	tail   *node
	act    *node
	length int
}

// ******** Public constructor ********

// New returns a new sorted value list.
func New() *SortedRuneList {
	return &SortedRuneList{head: nil, tail: nil, act: nil, length: 0}
}

// ******** Public functions ********

// Insert inserts a value with an associated index in the list.
func (l *SortedRuneList) Insert(index int, value rune) {
	// 1. Create the new node.
	newNode := &node{prev: nil, next: nil, value: value, index: index}

	// 2. If list is empty, start it with the new node.
	if l.head == nil {
		l.head = newNode
		l.tail = newNode
		l.length++
		return
	}

	// 3. List is not empty. Insert new node in list.
	insertNode(l.tail, newNode)

	// 4. Set new head or tail, if necessary.
	if newNode.prev == nil {
		l.head = newNode
	}

	if newNode.next == nil {
		l.tail = newNode
	}

	l.length++
}

// Length returns the length of the list.
func (l *SortedRuneList) Length() int {
	return l.length
}

// IsEmpty reports whether the list is empty.
func (l *SortedRuneList) IsEmpty() bool {
	return l.length == 0
}

// IsEndOfList reports whether the current list position is at the end of the list.
func (l *SortedRuneList) IsEndOfList() bool {
	return l.act == nil
}

// ResetPosition sets the list position to the begin of the list.
func (l *SortedRuneList) ResetPosition() {
	l.act = l.head
}

// GetNext returns the next element in the list and advances the position in the list.
func (l *SortedRuneList) GetNext() (int, rune) {
	result := *(l.act)

	l.act = result.next

	return result.index, result.value
}

// ValueOrderedIndices converts the list elements into a slice of the indices in the order of the values.
func (l *SortedRuneList) ValueOrderedIndices() []int {
	result := make([]int, l.length)

	var i int
	l.ResetPosition()
	for ri := 0; !l.IsEndOfList(); ri++ {
		i, _ = l.GetNext()
		result[ri] = i
	}

	return result
}

// ******** Private functions ********

// insertNode insert the new node beginning from the tail.
func insertNode(tail *node, new *node) {
	// Why from the last element? If there is a node that has the same value
	// the new element has to be put *after* that one.
	current := tail

	// Loop through list to see where to insert the new node.
	for {
		// 1. If the new value is greater or equal to the current value, append it to the current node.
		if new.value >= current.value {
			current.append(new)

			break
		}

		// 2. If there has been no place to insert the new node, yet and there is no
		//    previous node, prepend the new node to the current one.
		if current.prev == nil {
			current.prepend(new)

			break
		}

		// 3. Move forward, if no place found, yet and there is a previous node.
		current = current.prev
	}
}
