package json6

import (
	"reflect"
	"text/scanner"
)

func (dec *decoder) decodeArray() error {
	dec.tokens = append(dec.tokens, '[')
	prepVals := make([]reflect.Value, 0)

MAIN_LOOP:
	for {
		char := dec.s.Next()
		switch char {
		// empty value
		case ',':
			dec.tokens = append(dec.tokens, char)
			continue

		// end of array
		case ']':
			dec.tokens = append(dec.tokens, char)
			break MAIN_LOOP

		// can be a value, can be a whitespace
		default:
			// continue if char is whitespace or line terminator
			if isCharWhiteSpace(char) || isCharLineTerm(char) {
				continue
			} else if char == scanner.EOF {
				// if char is EOF, return error
				return errUnexpectedEOF(dec.s.Pos().Line, dec.s.Pos().Column)
			}

			// try to parse a value
			val := new(interface{})
			refVal := reflect.ValueOf(val)
			refVal = refVal.Elem()
			prepVals = append(prepVals, refVal)
			childDec := decoder{val: refVal, s: dec.s, state: array}
			if err := childDec.scan(true, char); err != nil {
				return err
			}

			dec.tokens = append(dec.tokens, childDec.tokens...)
			dec.tokens = append(dec.tokens, childDec.lastChar)

			// if last char is ']', array is complete
			if childDec.lastChar == ']' {
				break MAIN_LOOP
			}
		}
	}

	switch dec.val.Kind() {
	case reflect.Interface:
		slice := reflect.MakeSlice(reflect.TypeOf([]interface{}{}), 0, 0)
		slice = reflect.Append(slice, prepVals...)
		dec.val.Set(slice)

	case reflect.Slice:
		arrLen := len(prepVals)
		slice := reflect.MakeSlice(dec.val.Type(), arrLen, arrLen)
		dec.val.Set(slice)
		for i, v := range prepVals {
			// v is interface, so we need to call v.Elem()
			v = v.Elem()
			elm := dec.val.Index(i)
			if !setVal(elm, v) {
				return errMissMatchVal(string(dec.tokens)+" (inside array, type "+v.Type().String()+")", v.Type().Name(), v.Type().String())
			}
		}

	case reflect.Array:
		valLen := dec.val.Len()
		arrLen := len(prepVals)
		if arrLen > valLen {
			prepVals = prepVals[0:valLen]
		}

		for i, v := range prepVals {
			// v is interface, so we need to call v.Elem()
			v = v.Elem()
			elm := dec.val.Index(i)
			if !setVal(elm, v) {
				return errMissMatchVal(string(dec.tokens)+" (inside array, type "+v.Type().String()+")", v.Type().Name(), v.Type().String())
			}
		}

	default:
		return errMissMatchVal(string(dec.tokens)+" (array)", dec.val.Type().Name(), dec.val.Type().String())
	}

	return nil
}
