package json6

import (
	"errors"
	"fmt"
)

func errInvalidChar(invChar rune, pos *Position, nearChars []rune, expecting string) error {
	return fmt.Errorf("invalid character '%s' at %d:%d near '%s', expecting %s", string([]rune{invChar}), pos.Line(), pos.Column(), string(nearChars), expecting)
}

func errUnexpectedEOF(pos *Position, expecting string) error {
	return fmt.Errorf("unexpected EOF at %d:%d, expecting %s", pos.Line(), pos.Column(), expecting)
}

var ErrNoMoreToken = errors.New("no more token")
var ErrAlreadyAtBeginning = errors.New("already at beginning")
