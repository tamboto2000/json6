package json6

import "text/scanner"

var undefinedChars = []rune{'n', 'd', 'e', 'f', 'i', 'n', 'e', 'd'}

func decodeUndefined(s *scanner.Scanner) (*object, error) {
	for _, c := range undefinedChars {
		char := s.Next()
		if char != c {
			return nil, errInvalidChar(s.Pos().Line, s.Pos().Column, char, string(c))
		}
	}

	// check last char
	char := s.Peek()
	if isCharWhiteSpace(char) || isCharLineTerm(char) || isCharPunct(char) || char == scanner.EOF {
		return &object{kind: undefined}, nil
	}

	// advance reader if last char is invalid value terminator and return error
	s.Next()
	return nil, errInvalidChar(s.Pos().Line, s.Pos().Column, char, "whitespace, punctuator, or line terminator")
}
