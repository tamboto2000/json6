package json6

import "text/scanner"

var nullChars = []rune{'u', 'l', 'l'}

func decodeNull(s *scanner.Scanner) (*object, error) {
	for _, c := range nullChars {
		char := s.Next()
		if char != c {
			return nil, errInvalidChar(s.Pos().Line, s.Pos().Column, char, string(c))
		}
	}

	// check last char
	char := s.Peek()
	if isCharWhiteSpace(char) || isCharLineTerm(char) || isCharPunct(char) || char == scanner.EOF {
		return &object{kind: null}, nil
	}

	// advance reader if last char is invalid value terminator and return error
	s.Next()
	return nil, errInvalidChar(s.Pos().Line, s.Pos().Column, char, "whitespace, punctuator, or line terminator")
}
