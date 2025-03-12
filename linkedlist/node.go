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

package linkedlist

// ******** Private types ********

// node contains the data needed for a node in the doubly-linked list.
type node struct {
	prev  *node
	next  *node
	value rune
	index int
}

// ******** Private functions ********

// prepend prepends a new node to the current node.
func (c *node) prepend(n *node) {
	n.next = c
	n.prev = c.prev
	c.prev = n

	if n.prev != nil {
		n.prev.next = c.prev
	}
}

// append appends a node to the current node.
func (c *node) append(n *node) {
	n.prev = c
	n.next = c.next
	c.next = n

	if n.next != nil {
		n.next.prev = c.next
	}
}
