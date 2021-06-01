package json6

import (
	"bytes"
	"io"
	"strings"
	"text/scanner"
)

func scanFromBytes(byts []byte) (*object, error) {
	r := bytes.NewReader(byts)
	s := new(scanner.Scanner)
	s.Init(r)

	return scan(s)
}

func scanFromString(str string) (*object, error) {
	r := strings.NewReader(str)
	s := new(scanner.Scanner)
	s.Init(r)
	return scan(s)
}

func scanFromReader(r io.Reader) (*object, error) {
	s := new(scanner.Scanner)
	s.Init(r)
	return scan(s)
}

func scan(s *scanner.Scanner) (*object, error) {
	isEOF := false
	for !isEOF {
		char := s.Next()
		switch char {
		// undefined
		case 'u':
			return decodeUndefined(s)

		// null
		case 'n':
			return decodeNull(s)

		// booleans
		case 't':
			return decodeTrueBool(s)

		case 'f':
			return decodeFalseBool(s)

		case scanner.EOF:
			isEOF = true
			continue

		// if no any beginning of value is detected, check for whitespace and line terminator,
		// return error if character is neither of the two
		default:
			if !isCharWhiteSpace(char) && !isCharLineTerm(char) {
				return nil, errInvalidChar(s.Pos().Line, s.Pos().Column, char, "beginning of value")
			}
		}
	}

	return nil, nil
}
