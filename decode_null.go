package json6

import (
	"reflect"
	"text/scanner"
)

var nullChars = []rune{'u', 'l', 'l'}

func (dec *decoder) decodeNull() error {
	for _, c := range nullChars {
		char := dec.s.Next()
		if char != c {
			return errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, char, string(c))
		}
	}

	// check last char
	char := dec.s.Peek()
	if isCharWhiteSpace(char) || isCharLineTerm(char) || isCharPunct(char) || char == scanner.EOF {
		for {
			if dec.val.CanAddr() {
				dec.val = dec.val.Addr()
			} else {
				dec.val.Elem().Set(reflect.Zero(dec.val.Elem().Type()))
				break
			}
		}

		return nil
	}

	// advance reader if last char is invalid value terminator and return error
	dec.s.Next()
	return errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, char, "whitespace, punctuator, or line terminator")
}
