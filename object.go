package json6

import (
	"strconv"
	"text/scanner"
	"unicode"
)

func decodeObject(s *scanner.Scanner) (*object, error) {
	obj := &object{kind: obj, obj: make(map[string]*object)}

	for {
		isObjEnd := false
		for {
			char := s.Peek()
			if isCharWhiteSpace(char) || isCharLineTerm(char) {
				s.Next()
				continue
			} else if char == '}' {
				isObjEnd = true
				s.Next()
				break
			}

			break
		}

		if isObjEnd {
			break
		}

		// find an identifier name
		idName, err := _decodeIdName(s)
		if err != nil {
			return nil, err
		}

		// find punctuator ':'
		for {
			char := s.Next()
			if char == ':' {
				break
			} else if isCharWhiteSpace(char) || isCharLineTerm(char) {
				continue
			}

			return nil, errInvalidChar(s.Pos().Line, s.Pos().Column, char, "whitespace, line terminator, or ':'")
		}

		// find value
		val, err := scan(s)
		if err != nil {
			return nil, err
		}

		obj.obj[idName] = val

		// find punctuator ',' or '}'
		isEnd := false
		for {
			char := s.Next()
			if char == ',' {
				break
			} else if char == '}' {
				isEnd = true
				break
			} else if isCharWhiteSpace(char) || isCharLineTerm(char) {
				continue
			}

			return nil, errInvalidChar(s.Pos().Line, s.Pos().Column, char, "',', '}', or whitespace")
		}

		if isEnd {
			break
		}
	}

	return obj, nil
}

func _decodeIdName(s *scanner.Scanner) (string, error) {
	for {
		char := s.Next()
		switch char {
		case '`', '\'', '"':
			obj, err := decodeString(char, s)
			if err != nil {
				return "", err
			}

			return obj.str, nil

		case '_', '$', '\\':
			return _decodeUnquotedStr(char, s)

		default:
			if unicode.In(char, unicode.Lu, unicode.Ll, unicode.Lt, unicode.Lm, unicode.Lo, unicode.Nl) {
				return _decodeUnquotedStr(char, s)
			} else if isCharWhiteSpace(char) || isCharLineTerm(char) {
				continue
			}

			return "", errInvalidChar(s.Pos().Line, s.Pos().Column, char, "identifier start")
		}
	}
}

func _decodeUnquotedStr(begin rune, s *scanner.Scanner) (string, error) {
	rns := make([]rune, 0)
	// start with unicode
	if begin == '\\' {
		char := s.Next()
		if char != 'u' {
			return "", errInvalidChar(s.Pos().Line, s.Pos().Column, char, "'u'")
		}

		hexChars := make([]rune, 4)
		for i := 0; i < 4; i++ {
			char = s.Next()
			if !isCharValidHex(char) {
				return "", errInvalidChar(s.Pos().Line, s.Pos().Column, char, "hexadecimal number")
			}

			hexChars[i] = char
		}

		i, err := strconv.ParseInt(string(hexChars), 16, 32)
		if err != nil {
			return "", err
		}

		rns = append(rns, rune(i))
	} else {
		rns = append(rns, begin)
	}

	for {
		char := s.Peek()
		if char == '\\' {
			s.Next()
			// possible unicode
			char = s.Next()
			if char != 'u' {
				return "", errInvalidChar(s.Pos().Line, s.Pos().Column, char, "'u'")
			}

			hexChars := make([]rune, 4)
			for i := 0; i < 4; i++ {
				char = s.Next()
				if !isCharValidHex(char) {
					return "", errInvalidChar(s.Pos().Line, s.Pos().Column, char, "hexadecimal number")
				}

				hexChars[i] = char
			}

			i, err := strconv.ParseInt(string(hexChars), 16, 32)
			if err != nil {
				return "", err
			}

			rns = append(rns, rune(i))
			continue
		} else if isCharWhiteSpace(char) || isCharLineTerm(char) || isCharPunct(char) {
			break
		} else if char == scanner.EOF {
			s.Next()
			return "", errUnexpectedEOF(s.Pos().Line, s.Pos().Column)
		} else if !unicode.In(char,
			unicode.Lu,
			unicode.Ll,
			unicode.Lt,
			unicode.Lm,
			unicode.Lo,
			unicode.Nl,
			unicode.Mn,
			unicode.Mc,
			unicode.Nd,
			unicode.Pc) {
			s.Next()
			return "", errInvalidChar(s.Pos().Line, s.Pos().Column, char, "':', whitespace, or valid identifier character")
		}

		rns = append(rns, char)
		s.Next()
	}

	return string(rns), nil
}
