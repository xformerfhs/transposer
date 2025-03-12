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
