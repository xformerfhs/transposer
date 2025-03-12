package stringhelper

import "unicode"

// IsAlphaNumeric reports whether the supplied string contains only letters and digits.
func IsAlphaNumeric(s string) bool {
	for _, c := range s {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			return false
		}
	}

	return true
}
