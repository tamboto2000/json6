package json6

import "text/scanner"

func decodeComment(s *scanner.Scanner) (*object, error) {
	// check first char, if char is not '/' or '*', return error
	char := s.Next()
	// single line comment
	if char == '/' {
		for {
			if char = s.Next(); isCharLineTerm(char) || char == scanner.EOF {
				return &object{kind: comment}, nil
			}
		}
	} else if char == '*' {
		for {
			char = s.Next()
			switch char {
			// if another asterisk detected, check if it the end of comment
			case '*':
				char = s.Next()
				switch char {
				case '/':
					return &object{kind: comment}, nil

				case scanner.EOF:
					return nil, errUnexpectedEOF(s.Pos().Line, s.Pos().Column)
				}

			case scanner.EOF:
				return nil, errUnexpectedEOF(s.Pos().Line, s.Pos().Column)
			}
		}
	}

	return nil, errInvalidChar(s.Pos().Line, s.Pos().Column, char, "'/' or '*'")
}
