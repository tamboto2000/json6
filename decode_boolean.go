package json6

import (
	"reflect"
)

var (
	trueBoolChars  = []rune{'r', 'u', 'e'}
	falseBoolChars = []rune{'a', 'l', 's', 'e'}
)

func (dec *decoder) decodeTrueBool() error {
	for _, c := range trueBoolChars {
		char := dec.s.Next()
		if char != c {
			return dec.errInvalidChar(char, "'"+string(c)+"'")
		}
	}

	// check last char
	char := dec.s.Next()
	if err := dec.isCharEndOfValue(char); err != nil {
		return err
	}

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

func (dec *decoder) decodeFalseBool() error {
	for _, c := range falseBoolChars {
		char := dec.s.Next()
		if char != c {
			return dec.errInvalidChar(char, "'"+string(c)+"'")
		}
	}

	// check last char
	char := dec.s.Next()
	if err := dec.isCharEndOfValue(char); err != nil {
		return err
	}

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
