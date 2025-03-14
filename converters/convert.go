package converters

import "unicode"

type RuneConverter func(rune) rune

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
