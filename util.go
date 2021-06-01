package json6

import "unicode"

func isCharWhiteSpace(char rune) bool {
	return unicode.In(char, unicode.White_Space, unicode.Zs)
}

func isCharLineTerm(char rune) bool {
	switch char {
	case '\u000a', '\u000d', '\u2028', '\u2029':
		return true
	}

	return false
}

func isCharPunct(char rune) bool {
	switch char {
	case '{', '}', '[', ']', ':', ',':
		return true
	}

	return false
}

func isCharValidHex(char rune) bool {
	switch char {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true

	case 'a', 'b', 'c', 'd', 'e', 'f':
		return true

	case 'A', 'B', 'C', 'D', 'E', 'F':
		return true
	}

	return false
}
