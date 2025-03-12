package transposition

import (
	"slices"
	"transposer/slicehelper"
)

// UnTransposeRuneArray reverts a transposition with the given passwords.
func UnTransposeRuneArray(source []rune, passwords []string) []rune {
	sourceLength := len(source)
	result := make([]rune, sourceLength)

	from := result
	to := source
	// For decryption passwords have to be used last to first.
	slices.Reverse(passwords)
	for _, password := range passwords {
		from, to = to, from

		offsets := columnOrder(password)
		transposeLength := len(offsets)

		sourceIndex := 0
		for _, offset := range offsets {
			for destinationIndex := offset; destinationIndex < sourceLength; destinationIndex += transposeLength {
				to[destinationIndex] = from[sourceIndex]
				sourceIndex++
			}
		}

		slicehelper.ClearInteger(offsets)
	}

	return to
}
