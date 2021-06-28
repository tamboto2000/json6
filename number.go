package json6

import (
	"math"
	"strconv"
	"text/scanner"
	"unicode"
)

func decodeNumber(begin rune, s *scanner.Scanner) (*object, error) {
	obj := &object{kind: integer}

MAIN_SWITCH:
	switch begin {
	// handle underscore
	case '_':
		for {
			char := s.Next()
			if char == '_' {
				return nil, errInvalidChar(s.Pos().Line, s.Pos().Column, char, "'-', '+', 'I', 'N', '.', or numeric digit")
			}

			begin = char
			goto MAIN_SWITCH
		}

	// handle minus sign
	case '-':
		obj.numSign = '-'
		for {
			char := s.Next()
			switch char {
			case '-':
				if obj.numSign == '-' {
					obj.numSign = '+'
					continue
				} else {
					obj.numSign = '-'
					continue
				}

			default:
				begin = char
				goto MAIN_SWITCH
			}
		}

	// plus sign
	case '+':
		obj.numSign = '+'
		// check if after plus sign is valid char for number
		char := s.Peek()
		if !unicode.IsDigit(char) {
			s.Next()
			return nil, errInvalidChar(s.Pos().Line, s.Pos().Column, char, "numeric digit")
		}

		begin = s.Next()
		goto MAIN_SWITCH

	// Infinity
	case 'I':
		return obj, _decodeInfinity(obj, s)

	// NaN
	case 'N':
		return obj, _decodeNan(obj, s)

	// partial floating point
	case '.':
		return obj, _decodeFloatNumber(obj, s)

	// number start from 0
	case '0':
		char := s.Next()
		switch char {
		// exponent
		case 'e', 'E':
			obj.rns = append(obj.rns, []rune{'0', char}...)
			return obj, _decodeExpNumber(obj, s)

		// hexadecimal
		case 'x', 'X':
			return obj, _decodeHexNumber(obj, s)

		// octal
		case 'o':
			return obj, _decodeOctalNumber(obj, s)

		// binary
		case 'b':
			return obj, _decodeBinaryNumber(obj, s)

		default:
			return nil, errInvalidChar(s.Pos().Line, s.Pos().Column, char, "'.', 'e', 'E', 'x', 'X, 'o', or 'b'")
		}

	// number start from non-zero number
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		obj.rns = append(obj.rns, begin)
		for {
			char := s.Peek()
			if !unicode.IsDigit(char) {
				if isCharWhiteSpace(char) || isCharLineTerm(char) || isCharPunct(char) || char == scanner.EOF {
					break
				} else {
					switch char {
					case '_':
						s.Next()
						continue

					case '.':
						obj.rns = append(obj.rns, char)
						s.Next()
						return obj, _decodeFloatNumber(obj, s)

					case 'e', 'E':
						obj.rns = append(obj.rns, char)
						s.Next()
						return obj, _decodeExpNumber(obj, s)

					default:
						s.Next()
						return nil, errInvalidChar(s.Pos().Line, s.Pos().Column, char, "'_', '.', 'e', 'E', or numeric digit")
					}
				}
			}

			obj.rns = append(obj.rns, char)
			s.Next()
		}

		i, err := strconv.ParseInt(string(obj.rns), 10, 64)
		if err != nil {
			return nil, err
		}

		if obj.numSign == '-' {
			obj.integer = -i
		} else {
			obj.integer = i
		}

	default:
		return nil, errInvalidChar(s.Pos().Line, s.Pos().Column, begin, "beginning of number")
	}

	return obj, nil
}

var (
	infChars = []rune{'n', 'f', 'i', 'n', 'i', 't', 'y'}
	nanChars = []rune{'a', 'N'}
)

func _decodeInfinity(obj *object, s *scanner.Scanner) error {
	obj.kind = float
	for _, c := range infChars {
		char := s.Next()
		if char != c {
			return errInvalidChar(s.Pos().Line, s.Pos().Column, char, "'"+string(c)+"'")
		}

		obj.rns = append(obj.rns, char)
	}

	// check if last char is valid end of value
	char := s.Peek()
	if isCharWhiteSpace(char) || isCharLineTerm(char) || isCharPunct(char) || char == scanner.EOF {
		obj.kind = float
		if obj.numSign == '-' {
			obj.float = math.Inf(-1)
		} else {
			obj.float = math.Inf(1)
		}

		return nil
	}

	// advance reader if last char is invalid value terminator and return error
	s.Next()
	return errInvalidChar(s.Pos().Line, s.Pos().Column, char, "whitespace, punctuator, or line terminator")
}

func _decodeNan(obj *object, s *scanner.Scanner) error {
	obj.kind = float
	for _, c := range nanChars {
		char := s.Next()
		if char != c {
			return errInvalidChar(s.Pos().Line, s.Pos().Column, char, "'"+string(c)+"'")
		}

		obj.rns = append(obj.rns, char)
	}

	// check if last char is valid end of value
	char := s.Peek()
	if isCharWhiteSpace(char) || isCharLineTerm(char) || isCharPunct(char) || char == scanner.EOF {
		obj.kind = float
		if obj.numSign == '-' {
			obj.float = -math.NaN()
		} else {
			obj.float = math.NaN()
		}

		return nil
	}

	// advance reader if last char is invalid value terminator and return error
	s.Next()
	return errInvalidChar(s.Pos().Line, s.Pos().Column, char, "whitespace, punctuator, or line terminator")
}

func _decodeExpNumber(obj *object, s *scanner.Scanner) error {
	obj.kind = float

	// check first char
	char := s.Next()
	switch char {
	case '-', '+':
		obj.rns = append(obj.rns, char)
		char = s.Next()
		if !unicode.IsDigit(char) {
			return errInvalidChar(s.Pos().Line, s.Pos().Column, char, "numeric digit")
		}

		obj.rns = append(obj.rns, char)

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		obj.rns = append(obj.rns, char)

	default:
		return errInvalidChar(s.Pos().Line, s.Pos().Column, char, "'-', '+', or numeric digit")
	}

	for {
		char = s.Peek()
		if !unicode.IsDigit(char) {
			if isCharWhiteSpace(char) || isCharLineTerm(char) || isCharPunct(char) || char == scanner.EOF {
				break
			} else if char == '_' {
				s.Next()
				continue
			}

			return errInvalidChar(s.Pos().Line, s.Pos().Column, char, "whitespace, punctuator, line terminator, or numeric digit")
		}

		obj.rns = append(obj.rns, char)
		s.Next()
	}

	i, err := strconv.ParseFloat(string(obj.rns), 64)
	if err != nil {
		panic(err.Error())
	}

	if obj.numSign == '-' {
		obj.float = -i
	} else {
		obj.float = i
	}

	return nil
}

func _decodeHexNumber(obj *object, s *scanner.Scanner) error {
	for {
		char := s.Peek()
		if !isCharValidHex(char) {
			if isCharWhiteSpace(char) || isCharLineTerm(char) || isCharPunct(char) || char == scanner.EOF {
				break
			} else if char == '_' {
				s.Next()
				continue
			} else {
				s.Next()
				return errInvalidChar(s.Pos().Line, s.Pos().Column, char, "whitespace, punctuator, line terminator, or hexadecimal number")
			}
		}

		obj.rns = append(obj.rns, char)
		s.Next()
	}

	i, err := strconv.ParseInt(string(obj.rns), 16, 64)
	if err != nil {
		return err
	}

	if obj.numSign == '-' {
		obj.integer = -i
	} else {
		obj.integer = i
	}

	return nil
}

func _decodeOctalNumber(obj *object, s *scanner.Scanner) error {
	for {
		char := s.Peek()
		if !isCharValidOctal(char) {
			if isCharWhiteSpace(char) || isCharLineTerm(char) || isCharPunct(char) || char == scanner.EOF {
				break
			} else if char == '_' {
				s.Next()
				continue
			} else {
				s.Next()
				return errInvalidChar(s.Pos().Line, s.Pos().Column, char, "whitespace, punctuator, line terminator, or octal number")
			}
		}

		obj.rns = append(obj.rns, char)
		s.Next()
	}

	i, err := strconv.ParseInt(string(obj.rns), 8, 64)
	if err != nil {
		return err
	}

	if obj.numSign == '-' {
		obj.integer = -i
	} else {
		obj.integer = i
	}

	return nil
}

func _decodeBinaryNumber(obj *object, s *scanner.Scanner) error {
	for {
		char := s.Peek()
		if char != '0' && char != '1' {
			if isCharWhiteSpace(char) || isCharLineTerm(char) || isCharPunct(char) || char == scanner.EOF {
				break
			} else if char == '_' {
				s.Next()
				continue
			} else {
				s.Next()
				return errInvalidChar(s.Pos().Line, s.Pos().Column, char, "whitespace, punctuator, line terminator, or binary number")
			}
		}

		obj.rns = append(obj.rns, char)
		s.Next()
	}

	i, err := strconv.ParseInt(string(obj.rns), 2, 64)
	if err != nil {
		return err
	}

	if obj.numSign == '-' {
		obj.integer = -i
	} else {
		obj.integer = i
	}

	return nil
}

func _decodeFloatNumber(obj *object, s *scanner.Scanner) error {
	obj.kind = float
	for {
		char := s.Peek()
		if !unicode.IsDigit(char) {
			if isCharWhiteSpace(char) || isCharLineTerm(char) || isCharPunct(char) || char == scanner.EOF {
				break
			} else if char == '_' {
				s.Next()
				continue
			} else if char == 'e' || char == 'E' {
				obj.rns = append(obj.rns, char)
				s.Next()
				return _decodeExpNumber(obj, s)
			}

			s.Next()
			return errInvalidChar(s.Pos().Line, s.Pos().Column, char, "whitespace, punctuator, line terminator, numeric digit, or exponent")
		}

		obj.rns = append(obj.rns, char)
		s.Next()
	}

	i, err := strconv.ParseFloat(string(obj.rns), 64)
	if err != nil {
		return err
	}

	if obj.numSign == '-' {
		obj.float = -i
	} else {
		obj.float = i
	}

	return nil
}
