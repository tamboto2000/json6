package json6

import (	
	"reflect"
	"unicode"
)

func isCharWhiteSpace(char rune) bool {
	return unicode.In(char, unicode.White_Space, unicode.Zs)
}

const (
	lineFeed           = '\n'
	carriageReturn     = '\r'
	lineSeparator      = '\u2028'
	paragraphSeparator = '\u2029'
)

func isCharLineTerm(char rune) bool {
	switch char {
	case lineFeed, carriageReturn, lineSeparator, paragraphSeparator:
		return true
	}

	return false
}

func isCharPunct(char rune) bool {
	switch char {
	case '{', '}', '[', ']', ':', ',':
		return true
	}

	return false
}

func isCharValidHex(char rune) bool {
	switch char {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true

	case 'a', 'b', 'c', 'd', 'e', 'f':
		return true

	case 'A', 'B', 'C', 'D', 'E', 'F':
		return true
	}

	return false
}

func setVal(target, val reflect.Value) bool {
	switch target.Kind() {
	case reflect.Interface:
		target.Set(val)
		return true

	case reflect.Bool:
		if val.Kind() != reflect.Bool {
			return false
		}

		target.Set(val)
		return true

	case reflect.String:
		if val.Kind() != reflect.String {
			return false
		}

		target.Set(val)
		return true

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:		
		switch val.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			target.Set(val.Convert(target.Type()))
			return true

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			target.Set(val.Convert(target.Type()))
			return true

		default:			
			return false
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch val.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if val.Int() < 0 {
				return false
			}

			target.Set(val.Convert(target.Type()))
			return true

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			target.Set(val.Convert(target.Type()))
			return true

		default:
			return false
		}

	case reflect.Float32, reflect.Float64:
		switch val.Kind() {
		case reflect.Float32, reflect.Float64:
			target.Set(val.Convert(target.Type()))
			return true

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			target.Set(val.Convert(target.Type()))
			return true

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			target.Set(val.Convert(target.Type()))
			return true

		default:
			return false
		}

	case reflect.Slice:		
		switch val.Kind() {
		case reflect.Slice, reflect.Array:			
			arrLen := val.Len()
			slice := reflect.MakeSlice(target.Type(), arrLen, arrLen)
			target.Set(slice)
			for i := 0; i < arrLen; i++ {
				valElem := val.Index(i)
				if valElem.Kind() == reflect.Interface {
					valElem = valElem.Elem()
				}

				if !setVal(target.Index(i), valElem) {
					return false
				}
			}

			return true

		default:
			return false
		}

	case reflect.Array:
		switch val.Kind() {
		case reflect.Slice, reflect.Array:
			arrLen := val.Len()
			slice := reflect.MakeSlice(target.Type(), arrLen, arrLen)
			target.Set(slice)
			for i := 0; i < arrLen; i++ {
				valElem := val.Index(i)
				if valElem.Kind() == reflect.Interface {
					valElem = valElem.Elem()
				}

				if !setVal(target.Index(i), valElem) {
					return false
				}
			}

			return true

		default:
			return false
		}
	}

	return false
}
