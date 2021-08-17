package json6

import (
	"reflect"
)

var nullChars = []rune{'u', 'l', 'l'}

func (dec *decoder) decodeNull() error {
	for _, c := range nullChars {
		char := dec.s.Next()
		if char != c {
			return dec.errInvalidChar(char, string(c))
		}
	}

	// check last char
	char := dec.s.Next()
	if err := dec.isCharEndOfValue(char); err != nil {
		return err
	}

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
