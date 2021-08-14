package json6

import (
	"math"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

func (dec *decoder) decodeNumber(begin rune) error {
	dec.numSign = '+'
MAIN_SWITCH:
	switch begin {
	// handle minus sign (-)
	case '-':
		dec.numSign = '-'
		dec.tokens = append(dec.tokens, '-')
		for {
			char := dec.s.Next()
			if char == '-' {
				if dec.numSign == '+' {
					dec.numSign = '-'
				} else {
					dec.numSign = '+'
				}

				dec.tokens = append(dec.tokens, char)
			} else {
				if !isCharValidNumAfterSign(char) {
					return errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, char, "digits (0-9), '.', 'N', or 'I'")
				}

				begin = char
				break
			}
		}

		goto MAIN_SWITCH

	// handle plus sign (+)
	case '+':
		dec.tokens = append(dec.tokens, '+')
		char := dec.s.Next()
		if !isCharValidNumAfterSign(char) {
			return errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, char, "digits (0-9), '.', 'N', or 'I'")
		}

		begin = char
		goto MAIN_SWITCH

	// Infinity
	case 'I':
		return dec.decodeInfinity()

	// NaN
	case 'N':
		return dec.decodeNaN()

	// float num start with '.'
	case '.':
		dec.tokens = append(dec.tokens, begin)
		return dec.decodeFloatNum()

	// number start from 0
	case '0':
		dec.tokens = append(dec.tokens, begin)
		char := dec.s.Next()
		switch char {
		// digits
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			dec.tokens = append(dec.tokens, begin)

		// exponent
		case 'e', 'E':
			dec.tokens = append(dec.tokens, char)
			return dec.decodeExponent()

		// hexadecimal
		case 'x', 'X':
			dec.tokens = append(dec.tokens, char)
			return dec.decodeHexa()

		// binary
		case 'b', 'B':
			dec.tokens = append(dec.tokens, char)
			return dec.decodeBinary()

		// octal
		case 'o', 'O':
			dec.tokens = append(dec.tokens, char)
			return dec.decodeOctal()

		// floating point
		case '.':
			dec.tokens = append(dec.tokens, char)
			return dec.decodeFloatNum()
		}

	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		dec.tokens = append(dec.tokens, begin)
	}

LOOP:
	for {
		char := dec.s.Next()
		switch char {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			dec.tokens = append(dec.tokens, char)
			continue

		// exponent
		case 'e', 'E':
			dec.tokens = append(dec.tokens, char)
			return dec.decodeExponent()

		// float number
		case '.':
			dec.tokens = append(dec.tokens, char)
			return dec.decodeFloatNum()

		// separator
		case '_':
			dec.tokens = append(dec.tokens, char)
			if isEndOfVal, _, err := dec.validateNumUnderscore(); err != nil {
				return err
			} else {
				if isEndOfVal {
					break LOOP
				}
			}

			continue

		default:
			if !isCharEndOfValue(char) {
				return errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, char, "'_', '.', 'e', 'E', or digits (0-9)")
			}

			break LOOP
		}
	}

	i, err := strconv.ParseInt(prepNumTokensParseable(string(dec.tokens)), 10, 64)
	if err != nil {
		panic(err.Error())
	}

	if dec.numSign == '-' {
		i = -i
	}

	switch dec.val.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if i < 0 {
			return errMissMatchVal(string(dec.tokens)+" (int)", dec.val.Type().Name(), dec.val.Type().String())
		}

		dec.val.SetUint(uint64(i))

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		dec.val.SetInt(i)

	case reflect.Interface:
		dec.val.Set(reflect.ValueOf(i))

	default:
		return errMissMatchVal(string(dec.tokens)+" (int)", dec.val.Type().Name(), dec.val.Type().String())
	}

	return nil
}

var infChars = []rune{'n', 'f', 'i', 'n', 'i', 't', 'y'}

func (dec *decoder) decodeInfinity() error {
	for _, c := range infChars {
		if char := dec.s.Next(); char != c {
			return errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, char, string(c))
		}
	}

	// check last char
	char := dec.s.Peek()
	if isCharEndOfValue(char) {
		var sign int
		if dec.numSign == '-' {
			sign = -1
		} else {
			sign = +1
		}

		switch dec.val.Kind() {
		case reflect.Interface:
			dec.val.Set(reflect.ValueOf(math.Inf(sign)))

		case reflect.Float64:
			dec.val.SetFloat(math.Inf(sign))

		default:
			return errMissMatchVal("Infinity (float64)", dec.val.Type().Name(), dec.val.Type().String())
		}

		return nil
	}

	// advance reader if last char is invalid value terminator and return error
	dec.s.Next()
	return errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, char, "whitespace, punctuator, line terminator, or EOF")
}

var nanChars = []rune{'a', 'N'}

func (dec *decoder) decodeNaN() error {
	for _, c := range nanChars {
		if char := dec.s.Next(); char != c {
			return errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, char, string(c))
		}
	}

	// check last char
	char := dec.s.Peek()
	if isCharEndOfValue(char) {
		switch dec.val.Kind() {
		case reflect.Interface:
			dec.val.Set(reflect.ValueOf(math.NaN()))

		case reflect.Float64:
			dec.val.SetFloat(math.NaN())

		default:
			return errMissMatchVal("NaN (float64)", dec.val.Type().Name(), dec.val.Type().String())
		}

		return nil
	}

	// advance reader if last char is invalid value terminator and return error
	dec.s.Next()
	return errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, char, "whitespace, punctuator, line terminator, or EOF")
}

func (dec *decoder) decodeExponent() error {
	char := dec.s.Next()
	switch char {
	case '+', '-':
		dec.tokens = append(dec.tokens, char)
		char = dec.s.Next()
		if !unicode.IsNumber(char) {
			return errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, char, "digits (0-9)")
		}

		dec.tokens = append(dec.tokens, char)

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		dec.tokens = append(dec.tokens, char)

	default:
		return errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, char, "'-', '+', or digits (0-9)")
	}

LOOP:
	for {
		char = dec.s.Next()
		switch char {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			dec.tokens = append(dec.tokens, char)
			continue

		case '_':
			dec.tokens = append(dec.tokens, char)
			if isEndOfVal, _, err := dec.validateNumUnderscore(); err != nil {
				return err
			} else {
				if isEndOfVal {
					break LOOP
				}
			}

			continue

		default:
			if !isCharEndOfValue(char) {
				return errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, char, "'_', white space, line terminator, punctuation, EOF, or digits (0-9)")
			}

			break LOOP
		}
	}

	i, err := strconv.ParseFloat(prepNumTokensParseable(string(dec.tokens)), 64)
	if err != nil {
		panic(err.Error())
	}

	if dec.numSign == '-' {
		i = -i
	}

	switch dec.val.Kind() {
	case reflect.Float32, reflect.Float64:
		dec.val.SetFloat(i)

	case reflect.Interface:
		dec.val.Set(reflect.ValueOf(i))

	default:
		return errMissMatchVal(string(dec.tokens)+" (float)", dec.val.Type().Name(), dec.val.Type().String())
	}

	return nil
}

func (dec *decoder) decodeHexa() error {
LOOP:
	for {
		char := dec.s.Next()
		switch char {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			dec.tokens = append(dec.tokens, char)
			continue

		case 'a', 'b', 'c', 'd', 'e', 'f':
			dec.tokens = append(dec.tokens, char)
			continue

		case 'A', 'B', 'C', 'D', 'E', 'F':
			dec.tokens = append(dec.tokens, char)
			continue

		case '_':
			dec.tokens = append(dec.tokens, char)
			if isEndOfVal, _, err := dec.validateNumUnderscore(); err != nil {
				return err
			} else {
				if isEndOfVal {
					break LOOP
				}
			}

			continue

		default:
			if !isCharEndOfValue(char) {
				return errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, char, "'_', white space, line terminator, punctuation, EOF, or hexadecimal digits (0-9, a-f, A-F)")
			}

			break LOOP
		}
	}

	i, err := strconv.ParseInt(prepNumTokensParseable(string(dec.tokens)), 0, 32)
	if err != nil {
		panic(err.Error())
	}

	if dec.numSign == '-' {
		i = -i
	}

	switch dec.val.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if i < 0 {
			return errMissMatchVal(string(dec.tokens)+" (int)", dec.val.Type().Name(), dec.val.Type().String())
		}

		dec.val.SetUint(uint64(i))

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		dec.val.SetInt(i)

	case reflect.Interface:
		dec.val.Set(reflect.ValueOf(i))

	default:
		return errMissMatchVal(string(dec.tokens)+" (int)", dec.val.Type().Name(), dec.val.Type().String())
	}

	return nil
}

func (dec *decoder) decodeFloatNum() error {
	var lastChar rune
LOOP:
	for {
		char := dec.s.Next()
		switch char {
		case 'e', 'E':
			dec.tokens = append(dec.tokens, char)
			return dec.decodeExponent()

		case '_':
			// check last char in dec.tokens
			dec.tokens = append(dec.tokens, '_')
			if isEndOfval, lastC, err := dec.validateNumUnderscore(); err != nil {
				return err
			} else {
				if isEndOfval {
					lastChar = lastC
					break LOOP
				}
			}

			continue

		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			dec.tokens = append(dec.tokens, char)
			continue

		default:
			if isCharEndOfValue(char) {
				break LOOP
			}

			return errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, char, "'e', 'E', '_', or digits (0-9)")
		}
	}

	isParseable := false
	for _, token := range dec.tokens {
		if unicode.IsDigit(token) {
			isParseable = true
			break
		}
	}

	if !isParseable {
		return errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, lastChar, "digits (0-9)")
	}

	i, err := strconv.ParseFloat(prepNumTokensParseable(string(dec.tokens)), 64)
	if err != nil {
		panic(err.Error())
	}

	if dec.numSign == '-' {
		i = -i
	}

	switch dec.val.Kind() {
	case reflect.Float32, reflect.Float64:
		dec.val.SetFloat(i)

	case reflect.Interface:
		dec.val.Set(reflect.ValueOf(i))

	default:
		return errMissMatchVal(string(dec.tokens)+" (float)", dec.val.Type().Name(), dec.val.Type().String())
	}

	return nil
}

func (dec *decoder) decodeBinary() error {
LOOP:
	for {
		char := dec.s.Next()
		switch char {
		case '1', '0':
			dec.tokens = append(dec.tokens, char)
			continue

		case '_':
			dec.tokens = append(dec.tokens, char)
			if isEndOfVal, _, err := dec.validateNumUnderscore(); err != nil {
				return err
			} else {
				if isEndOfVal {
					break LOOP
				}
			}

			continue

		default:
			if isCharEndOfValue(char) {
				break LOOP
			}

			return errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, char, "'1', '0', or '_'")
		}
	}

	i, err := strconv.ParseInt(prepNumTokensParseable(string(dec.tokens)), 0, 64)
	if err != nil {
		panic(err.Error())
	}

	if dec.numSign == '-' {
		i = -i
	}

	switch dec.val.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if i < 0 {
			return errMissMatchVal(string(dec.tokens)+" (int)", dec.val.Type().Name(), dec.val.Type().String())
		}

		dec.val.SetUint(uint64(i))

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		dec.val.SetInt(i)

	case reflect.Interface:
		dec.val.Set(reflect.ValueOf(i))

	default:
		return errMissMatchVal(string(dec.tokens)+" (int)", dec.val.Type().Name(), dec.val.Type().String())
	}

	return nil
}

func (dec *decoder) decodeOctal() error {
LOOP:
	for {
		char := dec.s.Next()
		switch char {
		case '0', '1', '2', '3', '4', '5', '6', '7':
			dec.tokens = append(dec.tokens, char)
			continue

		case '_':
			dec.tokens = append(dec.tokens, char)
			if isEndOfVal, _, err := dec.validateNumUnderscore(); err != nil {
				return err
			} else {
				if isEndOfVal {
					break LOOP
				}
			}

			continue

		default:
			if !isCharEndOfValue(char) {
				return errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, char, "'_', white space, line terminator, punctuation, EOF, or octal digits (0-7)")
			}

			break LOOP
		}
	}

	i, err := strconv.ParseInt(prepNumTokensParseable(string(dec.tokens)), 0, 32)
	if err != nil {
		panic(err.Error())
	}

	if dec.numSign == '-' {
		i = -i
	}

	switch dec.val.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if i < 0 {
			return errMissMatchVal(string(dec.tokens)+" (int)", dec.val.Type().Name(), dec.val.Type().String())
		}

		dec.val.SetUint(uint64(i))

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		dec.val.SetInt(i)

	case reflect.Interface:
		dec.val.Set(reflect.ValueOf(i))

	default:
		return errMissMatchVal(string(dec.tokens)+" (int)", dec.val.Type().Name(), dec.val.Type().String())
	}

	return nil
}

func (dec *decoder) validateNumUnderscore() (isEndOfVal bool, lastChar rune, err error) {
	char := dec.tokens[len(dec.tokens)-1]
	switch char {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.', '_':
		char = dec.s.Next()
		// digits, _, end of value
		switch char {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			dec.tokens = append(dec.tokens, char)
			return false, char, nil

		case '_':
			dec.tokens = append(dec.tokens, char)
			return false, char, nil

		default:
			if !isCharEndOfValue(char) {
				return false, 0, errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, char, "'_', white space, line terminator, punctuation, EOF, or digits (0-9)")
			}

			return true, char, nil
		}

	default:
		return false, char, errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, char, "'e', 'E' or digits (0-9)")
	}
}

func isCharValidNumAfterSign(char rune) bool {
	switch char {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true

	case '.', 'N', 'I':
		return true
	}

	return false
}

func prepNumTokensParseable(tokens string) string {
	tokens = strings.TrimLeft(tokens, "-+")
	tokens = strings.ReplaceAll(tokens, "_", "")

	return tokens
}
