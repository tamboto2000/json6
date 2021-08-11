package json6

import (
	"reflect"
	"strconv"
	"text/scanner"
)

func (dec *decoder) decodeString(begin rune) error {
	rns := make([]rune, 0)
	charBegin := begin

	isEndOfStr := false
	for !isEndOfStr {
		char := dec.s.Next()

		switch char {
		// possible escaped character
		case '\\':
			char = dec.s.Next()

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

			case 'a':
				rns = append(rns, '\a')
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
					char = dec.s.Next()
					if !isCharValidHex(char) {
						return errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, char, "hexadecimal number")
					}

					hexChars[i] = char
				}

				// calculate the hexadecimal to decimal
				dec, err := strconv.ParseInt(string(hexChars), 16, 32)
				if err != nil {
					return err
				}

				rns = append(rns, rune(dec))
				continue

			// unicode
			case 'u':
				char = dec.s.Next()

				// if char is valid hexadecimal, then the unicode must contains 4 hexadecimal digit,
				// otherwise error will returned
				if isCharValidHex(char) {
					hexChars := make([]rune, 4)
					hexChars[0] = char
					for i := 0; i < 3; i++ {
						char = dec.s.Next()
						if !isCharValidHex(char) {
							return errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, char, "hexadecimal number")
						}

						hexChars[i+1] = char
					}

					// calculate the hexadecimal to decimal
					dec, err := strconv.ParseInt(string(hexChars), 16, 32)
					if err != nil {
						return err
					}

					rns = append(rns, rune(dec))
					continue
				} else if char == '{' {
					// if char is openbracket, then the hexadecimal digit can be more or less than 4
					hexChars := make([]rune, 0)
					for {
						char = dec.s.Next()
						if !isCharValidHex(char) {
							if char == '}' {
								break
							}

							return errInvalidChar(dec.s.Pos().Line, dec.s.Pos().Column, char, "hexadecimal number or '}'")
						}

						hexChars = append(hexChars, char)
					}

					// calculate the hexadecimal to decimal
					dec, err := strconv.ParseInt(string(hexChars), 16, 32)
					if err != nil {
						return err
					}

					rns = append(rns, rune(dec))
					continue
				}

				continue

			// line terminators
			case '\r':
				char = dec.s.Next()
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
			if charBegin == '`' {
				isEndOfStr = true
				continue
			}

			rns = append(rns, char)
			continue

		case '\'':
			if charBegin == '\'' {
				isEndOfStr = true
				continue
			}

			rns = append(rns, char)
			continue

		case '"':
			if charBegin == '"' {
				isEndOfStr = true
				continue
			}

			rns = append(rns, char)
			continue

		case scanner.EOF:
			return errUnexpectedEOF(dec.s.Pos().Line, dec.s.Pos().Column)
		}

		rns = append(rns, char)
	}

	if dec.val.Kind() == reflect.Interface {
		dec.val.Set(reflect.ValueOf(string(rns)))
	} else if dec.val.Kind() == reflect.String {
		dec.val.SetString(string(rns))
	} else {
		return errMissMatchVal("'"+string(rns)+"' (string)", dec.val.Type().Name(), dec.val.Type().String())
	}

	return nil
}
