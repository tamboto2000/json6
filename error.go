package json6

import (
	"errors"
	"strconv"
)

type ErrUnexpectedEOF error
type ErrInvalidChar error
type ErrUnmarshalNil error
type ErrUnmarshalNonPtr error
type ErrMissMatchVal error
type ErrEmptySource error

func errUnexpectedEOF(ln, col int) error {
	return ErrUnexpectedEOF(errors.New("unexpected EOF at " + strconv.Itoa(ln) + ":" + strconv.Itoa(col)))
}

func errInvalidChar(ln, col int, invChar rune, expect string) error {
	return ErrInvalidChar(errors.New("invalid character " + string(invChar) + " at " + strconv.Itoa(ln) + ":" + strconv.Itoa(col) + ", expecting " + expect))
}

func errUnmarshalNilVal() error {
	return ErrUnmarshalNil(errors.New("can not unmarshal to nil value"))
}

func errUnmarshalNonPtr() error {
	return ErrUnmarshalNonPtr(errors.New("can not unmarshal to non pointer value"))
}

func errMissMatchVal(srcType, valTypeName, valType string) error {
	return ErrMissMatchVal(errors.New("can not unmarshal " + srcType + " to " + valTypeName + " (" + valType + ")"))
}

func errEmptySource() error {
	return ErrEmptySource(errors.New("source is empty"))
}
