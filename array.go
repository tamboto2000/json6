package json6

import (
	"text/scanner"
)

func decodeArray(s *scanner.Scanner) (*object, error) {
	obj := &object{kind: array}

	for {
		char := s.Peek()
		if isCharWhiteSpace(char) || isCharLineTerm(char) || char == ',' {
			s.Next()
			continue
		} else if char == ']' {
			break
		} else {
			item, err := scan(s)
			if err != nil {
				return nil, err
			}

			obj.array = append(obj.array, item)

			for {
				char = s.Next()
				if isCharWhiteSpace(char) || isCharLineTerm(char) {
					continue
				} else if char == ',' {
					break
				} else if char == ']' {
					return obj, nil
				} else {
					return nil, errInvalidChar(s.Pos().Line, s.Pos().Column, char, "',' or ']'")
				}
			}

			continue
		}
	}

	return obj, nil
}
