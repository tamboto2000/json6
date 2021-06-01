package json6

import (
	"strconv"
	"text/scanner"
)

func decodeString(begin rune, s *scanner.Scanner) (*object, error) {
	obj := new(object)
	rns := make([]rune, 0)

	switch begin {
	case '`':
		obj.kind = backTickstr

	case '"':
		obj.kind = doubleQuoteStr

	case '\'':
		obj.kind = singleQuoteStr
	}

	isEndOfStr := false
	for !isEndOfStr {
		char := s.Next()

		switch char {
		// possible escaped character
		case '\\':
			char = s.Next()

			switch char {
			case '\'':
				rns = append(rns, '\'')
				continue

			case '"':
				rns = append(rns, '"')
				continue

			case '\\':
				rns = append(rns, '\\')
				continue

			case 'b':
				rns = append(rns, '\b')
				continue

			case 'f':
				rns = append(rns, '\f')
				continue

			case 'n':
				rns = append(rns, '\n')
				continue

			case 'r':
				rns = append(rns, '\r')
				continue

			case 't':
				rns = append(rns, '\t')
				continue

			case 'v':
				rns = append(rns, '\v')
				continue

			case '0':
				rns = append(rns, '\u0000')
				continue

			// hexadecimal code
			case 'x':
				hexChars := make([]rune, 2)
				for i := 0; i < 2; i++ {
					char = s.Next()
					if !isCharValidHex(char) {
						return nil, errInvalidChar(s.Pos().Line, s.Pos().Column, char, "valid hexadecimal number")
					}

					hexChars[i] = char
				}

				// calculate the hexadecimal to decimal
				dec, err := strconv.ParseInt(string(hexChars), 16, 64)
				if err != nil {
					return nil, err
				}

				rns = append(rns, rune(dec))
				continue

			// unicode
			case 'u':
				char = s.Next()

				// if char is valid hexadecimal, then the unicode must contains 4 hexadecimal digit,
				// otherwise error will returned
				if isCharValidHex(char) {
					hexChars := make([]rune, 4)
					hexChars[0] = char
					for i := 0; i < 3; i++ {
						char = s.Next()
						if !isCharValidHex(char) {
							return nil, errInvalidChar(s.Pos().Line, s.Pos().Column, char, "hexadecimal number")
						}

						hexChars[i+1] = char
					}

					// calculate the hexadecimal to decimal
					dec, err := strconv.ParseInt(string(hexChars), 16, 64)
					if err != nil {
						return nil, err
					}

					rns = append(rns, rune(dec))
					continue
				} else if char == '{' {
					// if char is openbracket, then the hexadecimal digit can be more than 4
					hexChars := make([]rune, 0)
					for {
						char = s.Next()
						if !isCharValidHex(char) {
							if char == '}' {
								break
							}

							return nil, errInvalidChar(s.Pos().Line, s.Pos().Column, char, "hexadecimal number or '}'")
						}

						hexChars = append(hexChars, char)
					}

					// calculate the hexadecimal to decimal
					dec, err := strconv.ParseInt(string(hexChars), 16, 64)
					if err != nil {
						return nil, err
					}

					rns = append(rns, rune(dec))
					continue
				}

			// line terminators
			case '\r':
				char = s.Next()
				if char == '\n' {
					continue
				}

				rns = append(rns, char)
				continue

			case '\n', lineSeparator, paragraphSeparator:
				continue

			default:
				rns = append(rns, char)
				continue
			}

		// end of string
		case '`':
			if obj.kind == backTickstr {
				isEndOfStr = true
				continue
			}

			rns = append(rns, char)
			continue

		case '\'':
			if obj.kind == singleQuoteStr {
				isEndOfStr = true
				continue
			}

			rns = append(rns, char)
			continue

		case '"':
			if obj.kind == doubleQuoteStr {
				isEndOfStr = true
				continue
			}

			rns = append(rns, char)
			continue

		case scanner.EOF:
			return nil, errUnexpectedEOF(s.Pos().Line, s.Pos().Column)
		}

		rns = append(rns, char)
	}

	obj.str = string(rns)
	return obj, nil
}
