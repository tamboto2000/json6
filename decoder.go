package json6

import (
	"bytes"
	"io"
	"reflect"
	"strconv"
)

// valueType define value type of a JSON6 value
type valueType uint

// value types
const (
	valueString valueType = iota
	valueInteger
	valueDouble
	valueObject
	valueArray
	valueNull
	valueBoolean
	valueUndefined
)

// value contains decoded value from token or sequence of tokens (like array and objects)
type value struct {
	t        valueType
	strVal   string
	intVal   int64
	floatVal float64
	boolVal  bool
	objVal   map[string]value // if t == ValueObject
	arrVal   []value          // if t == ValueArray
}

// expectState determine what token is to be expected
type expectState uint

const (
	// for decoding object
	expectIdent                      expectState = iota // identifier, string
	expectPunctColon                                    // ':'
	expectValue                                         // any JSON6 value
	expectPunctComaOrCloseCurlBrack                     // ',', '}'
	expectIdentOrPunctCloseCurlBrack                    // identifier, '}'

	// for decoding array
	expectValueOrPunctComaOrCloseBrack // any JSON6 value, ',', ']'
	expectPunctComaOrCloseBrack        // ',', ']'

	expectCommentOrEOF // comment, EOF
)

// decoder decode tokens into JSON6 value
type decoder struct {
	lx     *Lexer
	refVal reflect.Value
	val    value
}

// newDecoderFromBytes initiate new decoder from []byte
func newDecoderFromBytes(byts []byte, val interface{}) (*decoder, error) {
	r := newReader(bytes.NewReader(byts), newPosition(1, 0))
	lx := NewLexer(r)
	if err := lx.FetchTokens(); err != nil {
		return nil, err
	}

	// need validation
	refVal := reflect.ValueOf(val)

	return &decoder{
		lx:     lx,
		refVal: refVal,
	}, nil
}

// decodeValue decode any JSON6 value
func (dec *decoder) decodeValue() error {
	expect := expectValue
MAIN_LOOP:
	for {
		token, err := dec.lx.ReadToken()
		if err != nil {
			if err == ErrNoMoreToken {
				switch expect {
				case expectValue:
					return errUnexpectedEndOfTokenStream("any JSON6 value")

				case expectCommentOrEOF:
					break MAIN_LOOP
				}
			}

			return err
		}

		switch expect {
		case expectValue:
			switch token.t {
			case TokenString:
				dec.val, err = decodeString(token.runeReader)
				if err != nil {
					return err
				}

				expect = expectCommentOrEOF
				continue

			case TokenNumber:
				switch token.tokenNumSubType {
				case tokenNumInteger:
					dec.val, err = decodeIntNumber(token.runeReader)
					if err != nil {
						return err
					}

				case tokenNumDouble:
					dec.val, err = decodeIntNumber(token.runeReader)
					if err != nil {
						return err
					}
				}

				expect = expectCommentOrEOF
				continue

			case TokenNull:
				dec.val = value{t: valueNull}

				expect = expectCommentOrEOF
				continue

			case TokenBool:
				dec.val = decodeBool(token.String())

				expect = expectCommentOrEOF
				continue

			case TokenUndefined:
				dec.val = value{t: valueUndefined}

				expect = expectCommentOrEOF
				continue

			case TokenPunctuator:
				char := token.chars[0]

				if char == '{' {
					dec.val, err = decodeObject(dec.lx.tokenReader)
					if err != nil {
						return err
					}
				} else if char == '[' {
					dec.val, err = decodeArray(dec.lx.tokenReader)
					if err != nil {
						return err
					}
				} else {
					return errUnexpectedToken(token, "any JSON6 value")
				}

				expect = expectCommentOrEOF
				continue

			case TokenComment:
				continue
			}

			return errUnexpectedToken(token, "any JSON6 value")

		case expectCommentOrEOF:
			if token.t != TokenComment {
				return errUnexpectedToken(token, "EOF")
			}
		}
	}

	return nil
}

func decodeObject(r *tokenReader) (value, error) {
	expect := expectIdentOrPunctCloseCurlBrack
	var ident string
	val := value{t: valueObject, objVal: make(map[string]value)}

	for {
		token, err := r.ReadToken()
		if err != nil {
			if err == ErrNoMoreToken {
				switch expect {
				case expectIdent:
					return val, errUnexpectedEndOfTokenStream("identifier", "string")

				case expectPunctColon:
					return val, errUnexpectedEndOfTokenStream("':'")

				case expectValue:
					return val, errUnexpectedEndOfTokenStream("any JSON6 value")

				case expectPunctComaOrCloseCurlBrack:
					return val, errUnexpectedEndOfTokenStream("','", "'}'")

				case expectIdentOrPunctCloseCurlBrack:
					return val, errUnexpectedEndOfTokenStream("identifier", "string", "'}'")
				}
			}

			return val, err
		}

		// ignore comment
		if token.t == TokenComment {
			continue
		}

		switch expect {
		case expectIdent:
			switch token.t {
			case TokenIdentifier:
				ident, err = decodeIdentifier(token.runeReader)
				if err != nil {
					return val, err
				}

				expect = expectPunctColon
				continue

			case TokenString:
				val, err := decodeString(token.runeReader)
				if err != nil {
					return val, err
				}

				ident = val.strVal
				expect = expectPunctColon
				continue
			}

			return val, errUnexpectedToken(token, "identifier", "string")

		case expectPunctColon:
			if token.t == TokenPunctuator {
				if token.chars[0] == ':' {
					expect = expectValue

					continue
				}
			}

			return val, errUnexpectedToken(token, "':'")

		case expectValue:
			switch token.t {
			case TokenNumber:
				if token.tokenNumSubType == tokenNumInteger {
					decVal, err := decodeIntNumber(token.runeReader)
					if err != nil {
						return val, err
					}

					val.objVal[ident] = decVal
				} else {
					decVal, err := decodeDoubleNumber(token.runeReader)
					if err != nil {
						return val, err
					}

					val.objVal[ident] = decVal
				}

				expect = expectPunctComaOrCloseCurlBrack
				continue

			case TokenString:
				decVal, err := decodeString(token.runeReader)
				if err != nil {
					return val, err
				}

				val.objVal[ident] = decVal

				expect = expectPunctComaOrCloseCurlBrack
				continue

			case TokenNull:
				val.objVal[ident] = value{t: valueNull}

				expect = expectPunctComaOrCloseCurlBrack
				continue

			case TokenBool:
				val.objVal[ident] = decodeBool(token.String())

				expect = expectPunctComaOrCloseCurlBrack
				continue

			case TokenUndefined:
				val.objVal[ident] = value{t: valueUndefined}

				expect = expectPunctComaOrCloseCurlBrack
				continue

			// can be object or array
			case TokenPunctuator:
				char := token.chars[0]
				if char == '{' {
					decVal, err := decodeObject(r)
					if err != nil {
						return val, err
					}

					val.objVal[ident] = decVal

					expect = expectPunctComaOrCloseCurlBrack
					continue
				} else if char == '[' {
					decVal, err := decodeArray(r)
					if err != nil {
						return val, err
					}

					val.objVal[ident] = decVal

					expect = expectPunctComaOrCloseCurlBrack
					continue
				}
			}

			return val, errUnexpectedToken(token, "any JSON6 value")

		case expectPunctComaOrCloseCurlBrack:
			if token.t == TokenPunctuator {
				char := token.chars[0]
				if char == ',' {
					expect = expectIdentOrPunctCloseCurlBrack
					continue
				} else if char == '}' {
					return val, nil
				}

				return val, errUnexpectedToken(token, "','", "'}'")
			}

			return val, errUnexpectedToken(token, "','", "'}'")

		case expectIdentOrPunctCloseCurlBrack:
			switch token.t {
			case TokenIdentifier:
				ident, err = decodeIdentifier(token.runeReader)
				if err != nil {
					return val, err
				}

				expect = expectPunctColon
				continue

			case TokenString:
				decVal, err := decodeString(token.runeReader)
				if err != nil {
					return val, err
				}

				ident = decVal.strVal
				expect = expectPunctColon
				continue

			case TokenPunctuator:
				if token.chars[0] == '}' {
					return val, nil
				}

				return val, errUnexpectedToken(token, "'}'")
			}

			return val, errUnexpectedToken(token, "identifier", "string", "'}'")
		}
	}
}

func decodeArray(r *tokenReader) (value, error) {
	val := value{t: valueArray}
	expect := expectValueOrPunctComaOrCloseBrack

MAIN_LOOP:
	for {
		token, err := r.ReadToken()
		if err != nil {
			if err == ErrNoMoreToken {
				switch expect {
				case expectValueOrPunctComaOrCloseBrack:
					return val, errUnexpectedEndOfTokenStream("any JSON6 value", "','", "']'")

				case expectPunctComaOrCloseBrack:
					return val, errUnexpectedEndOfTokenStream("','", "']'")
				}
			}

			return val, err
		}

		// ignore comment
		if token.t == TokenComment {
			continue
		}

		switch expect {
		case expectValueOrPunctComaOrCloseBrack:
			switch token.t {
			case TokenString:
				decVal, err := decodeString(token.runeReader)
				if err != nil {
					return val, err
				}

				val.arrVal = append(val.arrVal, decVal)
				expect = expectPunctComaOrCloseBrack
				continue

			case TokenNumber:
				switch token.tokenNumSubType {
				case tokenNumInteger:
					decVal, err := decodeIntNumber(token.runeReader)
					if err != nil {
						return val, err
					}

					val.arrVal = append(val.arrVal, decVal)

				case tokenNumDouble:
					decVal, err := decodeDoubleNumber(token.runeReader)
					if err != nil {
						return val, err
					}

					val.arrVal = append(val.arrVal, decVal)
				}

				expect = expectPunctComaOrCloseBrack
				continue

			case TokenNull:
				val.arrVal = append(val.arrVal, value{t: valueNull})

				expect = expectPunctComaOrCloseBrack
				continue

			case TokenBool:
				decVal := decodeBool(token.String())
				val.arrVal = append(val.arrVal, decVal)

				expect = expectPunctComaOrCloseBrack
				continue

			case TokenUndefined:
				val.arrVal = append(val.arrVal, value{t: valueUndefined})

				expect = expectPunctComaOrCloseBrack
				continue

			case TokenPunctuator:
				char := token.chars[0]
				switch char {
				case '{':
					decVal, err := decodeObject(r)
					if err != nil {
						return val, err
					}

					val.arrVal = append(val.arrVal, decVal)
					expect = expectPunctComaOrCloseBrack
					continue

				case '[':
					decVal, err := decodeArray(r)
					if err != nil {
						return val, err
					}

					val.arrVal = append(val.arrVal, decVal)
					expect = expectPunctComaOrCloseBrack
					continue

				case ',':
					val.arrVal = append(val.arrVal, value{t: valueNull})
					continue

				case ']':
					break MAIN_LOOP
				}

				return val, errUnexpectedToken(token, "any JSON6 value", "','", "']'")
			}

		case expectPunctComaOrCloseBrack:
			if token.t != TokenPunctuator {
				return val, errUnexpectedToken(token, "','", "']'")
			}

			if token.chars[0] == ',' {
				val.arrVal = append(val.arrVal, value{t: valueNull})

				expect = expectValueOrPunctComaOrCloseBrack
				continue
			} else if token.chars[0] == ']' {
				break MAIN_LOOP
			}

			return val, errUnexpectedToken(token, "','", "']'")
		}
	}

	return val, nil
}

func decodeString(r *runeReader) (value, error) {
	var decVal []rune
	strBegin, _, _ := r.ReadRune()

	for {
		char, _, _ := r.ReadRune()
		if char == '\\' {
			char, _, _ := r.ReadRune()
			switch char {
			case '\\':
				decVal = append(decVal, char)
				continue

			case 'x':
				decChar := decodeHexaEscape(r)
				decVal = append(decVal, decChar)
				continue

			case 'u':
				decChar := decodeUnicodeEscape(r)
				decVal = append(decVal, decChar)
				continue

			case '\n':
				continue

			case '\r':
				char, _, _ := r.ReadRune()
				if char == '\n' {
					continue
				}

				decVal = append(decVal, char)
				continue

			case '\u2028':
				continue

			case '\u2029':
				continue

			case 'a':
				decVal = append(decVal, '\a')
				continue

			case 'b':
				decVal = append(decVal, '\b')
				continue

			case 'f':
				decVal = append(decVal, '\f')
				continue

			case 'n':
				decVal = append(decVal, '\n')
				continue

			case 'r':
				decVal = append(decVal, '\r')
				continue

			case 't':
				decVal = append(decVal, '\t')
				continue

			case 'v':
				decVal = append(decVal, '\v')
				continue

			case '0':
				decVal = append(decVal, '\u0000')
				continue
			}

			decVal = append(decVal, char)
			continue
		}

		if char == strBegin {
			break
		}

		decVal = append(decVal, char)
	}

	return value{t: valueString, strVal: string(decVal)}, nil
}

func decodeBool(str string) value {
	if str == "false" {
		return value{t: valueBoolean, boolVal: false}
	}

	return value{t: valueBoolean, boolVal: true}
}

func decodeIntNumber(r *runeReader) (value, error) {
	char, _, _ := r.ReadRune()
	if char == '-' {
		isMinus := true
		count := 1
		for {
			char, _, _ := r.ReadRune()
			if char == '-' {
				count++
				if isMinus {
					isMinus = false
				} else {
					isMinus = true
				}

				continue
			}

			break
		}

		i, err := strconv.ParseInt(string(r.chars[count:]), 0, 64)
		if err != nil {
			panic(err.Error())
		}

		if isMinus {
			i = -i
		}

		return value{t: valueInteger, intVal: i}, nil
	}

	i, err := strconv.ParseInt(string(r.chars), 0, 64)
	if err != nil {
		panic(err.Error())
	}

	return value{t: valueInteger, intVal: i}, nil
}

func decodeDoubleNumber(r *runeReader) (value, error) {
	char, _, _ := r.ReadRune()
	if char == '-' {
		isMinus := true
		count := 1
		for {
			char, _, _ := r.ReadRune()
			if char == '-' {
				count++
				if isMinus {
					isMinus = false
				} else {
					isMinus = true
				}

				continue
			}

			break
		}

		i, err := strconv.ParseFloat(string(r.chars[count:]), 64)
		if err != nil {
			panic(err.Error())
		}

		if isMinus {
			i = -i
		}

		return value{t: valueDouble, floatVal: i}, nil
	} else if char == '+' {
		i, err := strconv.ParseFloat(string(r.chars[1:]), 64)
		if err != nil {
			panic(err.Error())
		}

		return value{t: valueDouble, floatVal: i}, nil
	}

	i, err := strconv.ParseFloat(string(r.chars), 0)
	if err != nil {
		panic(err.Error())
	}

	return value{t: valueDouble, floatVal: i}, nil
}

func decodeIdentifier(r *runeReader) (string, error) {
	var decVal []rune

	for {
		char, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}

			return "", err
		}

		if char == '\\' {
			char, _, _ := r.ReadRune()
			switch char {
			case 'x':
				decChar := decodeHexaEscape(r)
				decVal = append(decVal, decChar)
				continue

			case 'u':
				decChar := decodeUnicodeEscape(r)
				decVal = append(decVal, decChar)
				continue
			}
		}

		decVal = append(decVal, char)
	}

	return string(decVal), nil
}

func decodeUnicodeEscape(r *runeReader) rune {
	var rns []rune
	char, _, _ := r.ReadRune()
	if char != '{' {
		rns = append(rns, char)
		for i := 0; i < 3; i++ {
			char, _, _ := r.ReadRune()
			rns = append(rns, char)
		}

		i, err := strconv.ParseInt(string(rns), 16, 32)
		if err != nil {
			panic(err.Error())
		}

		return rune(i)
	}

	for {
		char, _, _ := r.ReadRune()
		if char == '}' {
			break
		}

		rns = append(rns, char)
	}

	i, err := strconv.ParseInt(string(rns), 16, 32)
	if err != nil {
		panic(err.Error())
	}

	return rune(i)
}

func decodeHexaEscape(r *runeReader) rune {
	var rns []rune
	for i := 0; i < 2; i++ {
		char, _, _ := r.ReadRune()
		rns = append(rns, char)
	}

	i, err := strconv.ParseInt(string(rns), 16, 32)
	if err != nil {
		panic(err.Error())
	}

	return rune(i)
}
