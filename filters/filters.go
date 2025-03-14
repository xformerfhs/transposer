package filters

import "unicode"

type RuneFilter func(rune) bool

// Pass allows all characters.
func Pass(r rune) bool {
	return true
}

// OnlyLetters allows only letters.
func OnlyLetters(r rune) bool {
	return unicode.IsLetter(r)
}
