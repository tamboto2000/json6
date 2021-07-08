package json6

import (
	"reflect"
	"text/scanner"
)

func UnmarshalBytes(data []byte, v interface{}) error {
	s := scanFromBytes(data)
	val, err := valToReflect(v)
	if err != nil {
		return err
	}

	dec := newDecoder(val, s)

	return dec.scan()
}

func newDecoder(val reflect.Value, s *scanner.Scanner) *decoder {
	return &decoder{val: val, s: s}
}

type decoder struct {
	val reflect.Value
	s   *scanner.Scanner
}

func (dec *decoder) scan() error {
	isEOF := false
	for !isEOF {
		char := dec.s.Next()
		switch char {
		// // undefined
		// case 'u':
		// 	return decodeUndefined(s)

		// // null
		// case 'n':
		// 	return decodeNull(s)

		// booleans
		case 't':
			return dec.decodeTrueBool()

		case 'f':
			return dec.decodeFalseBool()

		// string
		case '`', '"', '\'':
			return dec.decodeString(char)

		// // comment
		// case '/':
		// 	return decodeComment(s)

		// // number
		// case '_', '-', '+', 'I', 'N', '.':
		// 	return decodeNumber(char, s)

		// case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		// 	return decodeNumber(char, s)

		// // array
		// case '[':
		// 	return decodeArray(s)

		// // object
		// case '{':
		// 	return decodeObject(s)

		case scanner.EOF:
			isEOF = true
			continue

		// if no any beginning of value is detected, check for whitespace and line terminator,
		// return error if character is neither of the two
		default:
			if !isCharWhiteSpace(char) && !isCharLineTerm(char) {
				return errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, char, "beginning of value")
			}
		}
	}

	return nil
}

func (dec *decoder) assignVal(data reflect.Value) error {
	return nil
}

// valToReflect convert v to reflect.Value and check if v is pointer.
// If v is not pointer, ErrUnmarshalNonPtr will be returned, if v is pointer but v is nil, ErrUnmarshalNilVal wil be returned
func valToReflect(v interface{}) (reflect.Value, error) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Interface {
		for {
			val = val.Elem()
			if val.Kind() != reflect.Interface {
				break
			}
		}
	}

	if val.Kind() != reflect.Ptr {
		return val, errUnmarshalNonPtr()
	}

	if val.IsNil() {
		return val, errUnmarshalNilVal()
	}

	val = indirect(val, false)

	return val, nil
}

// indirect walks down v allocating pointers as needed,
// until it gets to a non-pointer.
// If it encounters an Unmarshaler, indirect stops and returns that.
// If decodingNull is true, indirect stops at the first settable pointer so it
// can be set to nil.
func indirect(v reflect.Value, decodingNull bool) reflect.Value {
	// Issue #24153 indicates that it is generally not a guaranteed property
	// that you may round-trip a reflect.Value by calling Value.Addr().Elem()
	// and expect the value to still be settable for values derived from
	// unexported embedded struct fields.
	//
	// The logic below effectively does this when it first addresses the value
	// (to satisfy possible pointer methods) and continues to dereference
	// subsequent pointers as necessary.
	//
	// After the first round-trip, we set v back to the original value to
	// preserve the original RW flags contained in reflect.Value.
	v0 := v
	haveAddr := false

	// If v is a named type and is addressable,
	// start with its address, so that if the type has pointer methods,
	// we find them.
	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		haveAddr = true
		v = v.Addr()
	}
	for {
		// Load value from interface, but only if the result will be
		// usefully addressable.
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() && (!decodingNull || e.Elem().Kind() == reflect.Ptr) {
				haveAddr = false
				v = e
				continue
			}
		}

		if v.Kind() != reflect.Ptr {
			break
		}

		if decodingNull && v.CanSet() {
			break
		}

		// Prevent infinite loop if v is an interface pointing to its own address:
		//     var v interface{}
		//     v = &v
		if v.Elem().Kind() == reflect.Interface && v.Elem().Elem() == v {
			v = v.Elem()
			break
		}
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}

		if haveAddr {
			v = v0 // restore original value after round-trip Value.Addr().Elem()
			haveAddr = false
		} else {
			v = v.Elem()
		}
	}

	return v
}
