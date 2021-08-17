package json6

import (
	"text/scanner"
)

func (dec *decoder) decodeComment() error {
LOOP:
	for {
		char := dec.s.Next()
		switch char {
		// inline comment
		case '/':
			for {
				char = dec.s.Next()
				if isCharLineTerm(char) || char == scanner.EOF {
					break
				}
			}

			break LOOP

		// multi-line comment
		case '*':
		INNER_LOOP:
			for {
				char = dec.s.Next()
				switch char {
				case '*':
					char = dec.s.Next()
					switch char {
					case '/':
						break INNER_LOOP

					case scanner.EOF:
						return errUnexpectedEOF(dec.s.Pos().Line, dec.s.Pos().Column)
					}

				case scanner.EOF:
					return errUnexpectedEOF(dec.s.Pos().Line, dec.s.Pos().Column)
				}
			}

			break LOOP

		default:
			return dec.errInvalidChar(char, "'/' or '*'")
		}
	}

	return nil
}
