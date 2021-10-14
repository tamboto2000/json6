package json6

import (
	"bytes"
	"fmt"
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
	rnReader *runeReader
	strVal   string
	intVal   int64
	floatVal float64
	boolVal  bool
	objVal   map[string]value // if t == ValueObject
	arrVal   []value          // if t == ValueArray
}

// getVal get value based on value.t
func getVal(val value) reflect.Value {
	switch val.t {
	case valueString:
		return reflect.ValueOf(val.strVal)

	case valueInteger:
		return reflect.ValueOf(val.intVal)

	case valueDouble:
		return reflect.ValueOf(val.floatVal)

	case valueBoolean:
		return reflect.ValueOf(val.boolVal)

	case valueObject:
		m := make(map[string]interface{})
		for k, v := range val.objVal {
			m[k] = getVal(v).Interface()
		}

		return reflect.ValueOf(m)

	case valueArray:
		arr := make([]interface{}, 0)
		for _, v := range val.arrVal {
			arr = append(arr, getVal(v).Interface())
		}

		return reflect.ValueOf(arr)
	}

	return reflect.ValueOf(nil)
}

// getValTypeStr get value type in string
func getValTypeStr(val value) string {
	switch val.t {
	case valueString:
		return "string"

	case valueInteger:
		return "integer"

	case valueDouble:
		return "double"

	case valueBoolean:
		return "boolean"

	case valueObject:
		return "object"

	case valueArray:
		return "array"

	case valueNull:
		return "null"

	case valueUndefined:
		return "undefined"
	}

	return ""
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

	refVal, err := valToReflect(val)
	if err != nil {
		return nil, err
	}

	return &decoder{
		lx:     lx,
		refVal: refVal,
	}, nil
}

func valToReflect(v interface{}) (reflect.Value, error) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Interface {
		for {
			val = val.Elem()
			if val.Kind() != reflect.Interface {
				break
			}
		}
	}

	if val.Kind() != reflect.Ptr {
		return val, errDecodeToNonPtr()
	}

	if val.IsNil() {
		return val, errDecodeToNilPtr()
	}

	val = indirect(val, false)

	return val, nil
}

func assignValue(refVal reflect.Value, val *value) error {
	switch val.t {
	case valueObject:
		return assignObjectValue(refVal, val)
	case valueArray:
		return assignArrayValue(refVal, val)
	case valueString:
		return assignStrValue(refVal, val)

	case valueInteger:
		return assignIntNumValue(refVal, val)

	case valueDouble:
		return assignDoubleNumValue(refVal, val)

	case valueNull:
		return assignNullValue(refVal, val)

	case valueBoolean:
		return assignBoolValue(refVal, val)

	case valueUndefined:
		return assignUndefinedValue(refVal, val)
	}

	return nil
}

func assignObjectValue(refVal reflect.Value, val *value) error {
	switch refVal.Kind() {
	case reflect.Struct:
		// temporary storage for field
		// accepted tags are json6, json5, and json
		storeFields := make(map[string]reflect.Value)
		storeFieldNames := make(map[string]string)
		numField := refVal.NumField()
		for i := 0; i < numField; i++ {
			field := refVal.Type().Field(i)
			if tag := field.Tag.Get("json6"); tag != "" {
				storeFields[tag] = refVal.Field(i)
				storeFieldNames[tag] = field.Name
			} else if tag := field.Tag.Get("json5"); tag != "" {
				storeFields[tag] = refVal.Field(i)
				storeFieldNames[tag] = field.Name
			} else if tag := field.Tag.Get("json"); tag != "" {
				storeFields[tag] = refVal.Field(i)
				storeFieldNames[tag] = field.Name
			} else {
				storeFields[field.Name] = refVal.Field(i)
				storeFieldNames[field.Name] = field.Name
			}
		}

		for k, v := range val.objVal {
			if rv, ok := storeFields[k]; ok {
				vc := v
				err := assignValue(rv, &vc)
				if err != nil {
					return fmt.Errorf("error decoding value to %s.%s:\n%s", refVal.Type().Name(), storeFieldNames[k], err.Error())
				}
			}
		}

	case reflect.Map:
		// verify if the map key type is string
		refType := refVal.Type()
		if refType.Key().Kind() != reflect.String {
			return fmt.Errorf("can not decode object to map[%s]%s, JSON6 object can only be decoded to struct or map[string]interface{}",
				refType.Key().String(), refType.Elem().String())
		}

		// verify if the map elem type is interface{}
		if refType.Elem().Kind() != reflect.Interface {
			return fmt.Errorf("can not decode object to map[%s]%s, JSON6 object can only be decoded to struct or map[string]interface{}",
				refType.Key().String(), refType.Elem().String())
		}

		for k, v := range val.objVal {
			refVal.SetMapIndex(reflect.ValueOf(k), getVal(v))
		}

	case reflect.Interface:
		refVal.Set(getVal(*val))

	default:
		return fmt.Errorf("can not decode object to %s (%s), JSON6 object can only be decoded to struct or map[string]interface{}", refVal.Type().Name(), refVal.Type().String())
	}

	return nil
}

func assignArrayValue(refVal reflect.Value, val *value) error {
	switch refVal.Kind() {
	case reflect.Slice:
		refValElemKindStr := refVal.Type().Elem().String()
		refValElemType := refVal.Type().Elem()
		refVal.SetLen(0)

		for _, v := range val.arrVal {
			if v.t == valueNull || v.t == valueUndefined {
				zeroVal := reflect.Zero(refValElemType)
				refVal.Set(reflect.Append(refVal, zeroVal))
				continue
			}

			refV := getVal(v)

			if !refV.Type().ConvertibleTo(refValElemType) {
				return errMismatchType(v.rnReader.chars, getValTypeStr(v), refValElemKindStr)
			}

			refVal.Set(reflect.Append(refVal, refV.Convert(refValElemType)))
		}

	case reflect.Array:
		refValElemType := refVal.Type().Elem()
		for i := 0; i < refVal.Len(); i++ {
			v := val.arrVal[i]
			if v.t == valueNull || v.t == valueUndefined {
				zeroVal := reflect.Zero(refValElemType)
				refVal.Index(i).Set(zeroVal)
				continue
			}

			if err := assignValue(refVal.Index(i), &v); err != nil {
				return err
			}
		}

	case reflect.Interface:
		var arr []interface{}
		for _, v := range val.arrVal {
			if v.t == valueNull || v.t == valueUndefined {
				arr = append(arr, nil)

				continue
			}

			arr = append(arr, getVal(v).Interface())
		}

		refVal.Set(reflect.ValueOf(arr))

	default:
		return errMismatchType(val.rnReader.chars, "array", refVal.Type().String())
	}

	return nil
}

func assignStrValue(refVal reflect.Value, val *value) error {
	switch refVal.Kind() {
	case reflect.String:
		refVal.SetString(val.strVal)

	case reflect.Interface:
		refVal.Set(reflect.ValueOf(val.strVal))

	default:
		return errMismatchType(val.rnReader.chars, "string", refVal.Type().String())
	}

	return nil
}

func assignIntNumValue(refVal reflect.Value, val *value) error {
	switch refVal.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		refVal.SetInt(val.intVal)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if val.intVal < 0 {
			return errMismatchType(val.rnReader.chars, "integer", refVal.Type().String())
		}

		refVal.SetUint(uint64(val.intVal))

	case reflect.Interface:
		refVal.Set(reflect.ValueOf(val.intVal))

	default:
		return errMismatchType(val.rnReader.chars, "integer", refVal.Type().String())
	}

	return nil
}

func assignDoubleNumValue(refVal reflect.Value, val *value) error {
	switch refVal.Kind() {
	case reflect.Float32, reflect.Float64:
		refVal.SetFloat(val.floatVal)

	case reflect.Interface:
		refVal.Set(reflect.ValueOf(val.floatVal))

	default:
		return errMismatchType(val.rnReader.chars, "integer", refVal.Type().String())
	}

	return nil
}

func assignNullValue(refVal reflect.Value, val *value) error {
	if refVal.IsValid() {
		if !refVal.IsZero() {
			refVal.Set(reflect.Zero(refVal.Type()))
		}
	}

	return nil
}

func assignBoolValue(refVal reflect.Value, val *value) error {
	switch refVal.Kind() {
	case reflect.Bool:
		refVal.SetBool(val.boolVal)

	case reflect.Interface:
		refVal.Set(reflect.ValueOf(val.boolVal))

	default:
		return errMismatchType(val.rnReader.chars, "bool", refVal.Type().String())
	}

	return nil
}

func assignUndefinedValue(refVal reflect.Value, val *value) error {
	if refVal.IsValid() {
		if !refVal.IsZero() {
			refVal.Set(reflect.Zero(refVal.Type()))
		}
	}

	return nil
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
				dec.val = decodeBool(token.runeReader)

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

	return assignValue(dec.refVal, &dec.val)
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
				val.objVal[ident] = decodeBool(token.runeReader)

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
				decVal := decodeBool(token.runeReader)
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

	return value{t: valueString, strVal: string(decVal), rnReader: r}, nil
}

func decodeBool(r *runeReader) value {
	if string(r.chars) == "false" {
		return value{t: valueBoolean, boolVal: false}
	}

	return value{t: valueBoolean, boolVal: true, rnReader: r}
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

	return value{t: valueInteger, intVal: i, rnReader: r}, nil
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

	return value{t: valueDouble, floatVal: i, rnReader: r}, nil
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

// indirect walks down v allocating pointers as needed,
// until it gets to a non-pointer.
// Btw, this piece of code is stolen from encoding/json lol
func indirect(v reflect.Value, decodingNull bool) reflect.Value {
	// Issue #24153 indicates that it is generally not a guaranteed property
	// that you may round-trip a reflect.Value by calling Value.Addr().Elem()
	// and expect the value to still be settable for values derived from
	// unexported embedded struct fields.
	//
	// The logic below effectively does this when it first addresses the value
	// (to satisfy possible pointer methods) and continues to dereference
	// subsequent pointers as necessary.
	//
	// After the first round-trip, we set v back to the original value to
	// preserve the original RW flags contained in reflect.Value.
	v0 := v
	haveAddr := false

	// If v is a named type and is addressable,
	// start with its address, so that if the type has pointer methods,
	// we find them.
	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		haveAddr = true
		v = v.Addr()
	}

	for {
		// Load value from interface, but only if the result will be
		// usefully addressable.
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() && (!decodingNull || e.Elem().Kind() == reflect.Ptr) {
				haveAddr = false
				v = e
				continue
			}
		}

		if v.Kind() != reflect.Ptr {
			break
		}

		if decodingNull && v.CanSet() {
			break
		}

		// Prevent infinite loop if v is an interface pointing to its own address:
		//     var v interface{}
		//     v = &v
		if v.Elem().Kind() == reflect.Interface && v.Elem().Elem() == v {
			v = v.Elem()
			break
		}
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}

		if haveAddr {
			v = v0 // restore original value after round-trip Value.Addr().Elem()
			haveAddr = false
		} else {
			v = v.Elem()
		}
	}

	return v
}

func Unmarshal(src []byte, val interface{}) error {
	dec, err := newDecoderFromBytes(src, val)
	if err != nil {
		return err
	}

	return dec.decodeValue()
}
