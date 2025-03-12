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
