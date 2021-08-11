package json6

import (
	"math"
	"reflect"
	"strconv"
	"text/scanner"
	"unicode"
)

func (dec *decoder) decodeNumber(begin rune) error {
	dec.numSign = '+'
MAIN_SWITCH:
	switch begin {
	// handle minus sign (-)
	case '-':
		dec.numSign = '-'
		for {
			char := dec.s.Next()
			if char == '-' {
				if dec.numSign == '+' {
					dec.numSign = '-'
				} else {
					dec.numSign = '+'
				}
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

		// number start from 0
	case '0':
		char := dec.s.Next()
		switch char {

		// exponent
		case 'e', 'E':
			dec.tokens = append(dec.tokens, []rune{'0', char}...)
			return dec.decodeExponent()

		// hexadecimal
		case 'x', 'X':

		// floating point
		case '.':
			dec.tokens = append(dec.tokens, []rune{'0', char}...)
			return dec.decodeFloatNum()
		}
	}

	return nil
}

var infChars = []rune{'n', 'f', 'i', 'n', 'i', 't', 'y'}

func (dec *decoder) decodeInfinity() error {
	if dec.val.Kind() != reflect.Float64 && dec.val.Kind() != reflect.Interface {
		return errMissMatchVal("Infinity (float64)", dec.val.Type().Name(), dec.val.Type().String())
	}

	for _, c := range infChars {
		if char := dec.s.Next(); char != c {
			return errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, char, string(c))
		}
	}

	// check last char
	char := dec.s.Peek()
	if isCharWhiteSpace(char) || isCharLineTerm(char) || isCharPunct(char) || char == scanner.EOF {
		var sign int
		if dec.numSign == '-' {
			sign = -1
		} else {
			sign = +1
		}

		if dec.val.Kind() == reflect.Interface {
			dec.val.Set(reflect.ValueOf(math.Inf(sign)))
		} else {
			dec.val.SetFloat(math.Inf(sign))
		}

		return nil
	}

	// advance reader if last char is invalid value terminator and return error
	dec.s.Next()
	return errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, char, "whitespace, punctuator, or line terminator")
}

var nanChars = []rune{'a', 'N'}

func (dec *decoder) decodeNaN() error {
	if dec.val.Kind() != reflect.Float64 && dec.val.Kind() != reflect.Interface {
		return errMissMatchVal("NaN (float64)", dec.val.Type().Name(), dec.val.Type().String())
	}

	for _, c := range nanChars {
		if char := dec.s.Next(); char != c {
			return errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, char, string(c))
		}
	}

	// check last char
	char := dec.s.Peek()
	if isCharWhiteSpace(char) || isCharLineTerm(char) || isCharPunct(char) || char == scanner.EOF {
		if dec.val.Kind() == reflect.Interface {
			dec.val.Set(reflect.ValueOf(math.NaN()))
		} else {
			dec.val.SetFloat(math.NaN())
		}

		return nil
	}

	// advance reader if last char is invalid value terminator and return error
	dec.s.Next()
	return errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, char, "whitespace, punctuator, or line terminator")
}

func (dec *decoder) decodeExponent() error {
	if dec.val.Kind() != reflect.Float64 && dec.val.Kind() != reflect.Interface {
		return errMissMatchVal("float64", dec.val.Type().Name(), dec.val.Type().String())
	}

	// check first char
	char := dec.s.Next()
	switch char {
	case '-', '+':
		dec.tokens = append(dec.tokens, char)
		char = dec.s.Next()
		if !unicode.IsDigit(char) {
			return errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, char, "digits")
		}

		dec.tokens = append(dec.tokens, char)

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		dec.tokens = append(dec.tokens, char)

	default:
		return errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, char, "'-', '+', or digits")
	}

	for {
		char = dec.s.Peek()
		if !unicode.IsDigit(char) {
			if isCharWhiteSpace(char) || isCharLineTerm(char) || isCharPunct(char) || char == scanner.EOF {
				break
			} else if char == '_' {
				dec.s.Next()
				continue
			}

			dec.s.Next()
			return errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, char, "whitespace, punctuator, line terminator, or numeric digit")
		}

		dec.tokens = append(dec.tokens, char)
		dec.s.Next()
	}

	i, err := strconv.ParseFloat(string(dec.tokens), 64)
	if err != nil {
		panic(err.Error())
	}

	if dec.numSign == '-' {
		i = -i
	}

	if dec.val.Kind() == reflect.Interface {
		dec.val.Set(reflect.ValueOf(i))
	} else {
		dec.val.SetFloat(i)
	}

	return nil
}

func (dec *decoder) decodeHexa() error {
	

	return nil
}

func (dec *decoder) decodeFloatNum() error {
	if dec.val.Kind() != reflect.Float64 && dec.val.Kind() != reflect.Float32 && dec.val.Kind() != reflect.Interface {
		return errMissMatchVal("float64", dec.val.Type().Name(), dec.val.Type().String())
	}

	for {
		char := dec.s.Peek()
		if !unicode.IsDigit(char) {
			if isCharWhiteSpace(char) || isCharLineTerm(char) || isCharPunct(char) || char == scanner.EOF {
				break
			} else if char == '_' {
				dec.s.Next()
				continue
			} else if char == 'e' || char == 'E' {
				dec.tokens = append(dec.tokens, char)
				dec.s.Next()
				return dec.decodeExponent()
			}

			dec.s.Next()
			return errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, char, "whitespace, punctuator, line terminator, numeric digit, or exponent")
		}

		dec.tokens = append(dec.tokens, char)
		dec.s.Next()
	}

	i, err := strconv.ParseFloat(string(dec.tokens), 64)
	if err != nil {
		return err
	}

	if dec.numSign == '-' {
		i = -i
	}

	if dec.val.Kind() == reflect.Interface {
		dec.val.Set(reflect.ValueOf(i))
	} else {
		dec.val.SetFloat(i)
	}

	return nil
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

func isNumSignedInt(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	}

	return false
}

func isNumUnsignedInt(kind reflect.Kind) bool {
	switch kind {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	}

	return false
}
