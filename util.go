package json6

import (
	"text/scanner"
	"unicode"
)

func isCharWhiteSpace(char rune) bool {
	return unicode.In(char, unicode.White_Space, unicode.Zs)
}

const (
	lineFeed           = '\u000a'
	carriageReturn     = '\u000d'
	lineSeparator      = '\u2028'
	paragraphSeparator = '\u2029'
)

func isCharLineTerm(char rune) bool {
	switch char {
	case lineFeed, carriageReturn, lineSeparator, paragraphSeparator:
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

func isCharEndOfValue(char rune) bool {
	if isCharWhiteSpace(char) || isCharLineTerm(char) || isCharPunct(char) || char == scanner.EOF {
		return true
	}

	return false
}
