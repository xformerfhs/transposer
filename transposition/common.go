package transposition

import "transposer/linkedlist"

// setupResult creates the slices and lengths needed for transposing and
// untransposing.
func setupResult(source []rune) (int, []rune) {
	sourceLength := len(source)
	result := make([]rune, sourceLength)

	return sourceLength, result
}

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
