package json6

import (
	"errors"
	"fmt"
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

func (dec *decoder) errInvalidChar(invChar rune, expect string) error {
	return fmt.Errorf("invalid character '%s' at %d:%d, expecting %s", string([]rune{invChar}), dec.s.Pos().Line, dec.s.Pos().Column, expect)
}

func errUnmarshalNilVal() error {
	return errors.New("can not unmarshal to nil value")
}

func errUnmarshalNonPtr() error {
	return errors.New("can not unmarshal to non pointer value")
}

func errMissMatchVal(srcType, valTypeName, valType string) error {
	return errors.New("can not unmarshal " + srcType + " to " + valTypeName + " (" + valType + ")")
}

func errEmptySource() error {
	return errors.New("source is empty")
}
