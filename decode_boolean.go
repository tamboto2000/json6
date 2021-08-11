package json6

import (
	"reflect"
	"text/scanner"
)

var (
	trueBoolChars  = []rune{'r', 'u', 'e'}
	falseBoolChars = []rune{'a', 'l', 's', 'e'}
)

func (dec *decoder) decodeTrueBool() error {
	for _, c := range trueBoolChars {
		char := dec.s.Next()
		if char != c {
			return errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, char, "'"+string(c)+"'")
		}
	}

	// check last char
	char := dec.s.Peek()
	if isCharWhiteSpace(char) || isCharLineTerm(char) || isCharPunct(char) || char == scanner.EOF {
		switch dec.val.Kind() {
		case reflect.Interface:
			dec.val.Set(reflect.ValueOf(true))

		case reflect.Bool:
			dec.val.SetBool(true)

		default:
			return errMissMatchVal("true (boolean)", dec.val.Type().Name(), dec.val.Type().String())
		}

		return nil
	}

	// advance reader if last char is invalid value terminator and return error
	dec.s.Next()
	return errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, char, "whitespace, punctuator, or line terminator")
}

func (dec *decoder) decodeFalseBool() error {
	for _, c := range falseBoolChars {
		char := dec.s.Next()
		if char != c {
			return errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, char, "'"+string(c)+"'")
		}
	}

	// check last char
	char := dec.s.Peek()
	if isCharWhiteSpace(char) || isCharLineTerm(char) || isCharPunct(char) || char == scanner.EOF {
		switch dec.val.Kind() {
		case reflect.Interface:
			dec.val.Set(reflect.ValueOf(false))

		case reflect.Bool:
			dec.val.SetBool(false)

		default:
			return errMissMatchVal("false (boolean)", dec.val.Type().Name(), dec.val.Type().String())
		}

		return nil
	}

	// advance reader if last char is invalid value terminator and return error
	dec.s.Next()
	return errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, char, "whitespace, punctuator, or line terminator")
}
