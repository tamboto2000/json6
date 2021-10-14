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

func errUnexpectedToken(token Token, expects ...string) error {
	expectsLen := len(expects)
	var expectStr string
	if expectsLen > 1 {
		for i, expect := range expects {
			if i < expectsLen-1 {
				expectStr += expect + ", "
			} else {
				expectStr += "or " + expect
			}
		}
	} else {
		expectStr = expects[0]
	}

	return fmt.Errorf("unexpected token '%s' (%s) at %d:%d, expecting %s", token.String(), token.TypeString(), token.StartPos.ln, token.StartPos.col, expectStr)
}

func errUnexpectedEndOfTokenStream(expects ...string) error {
	expectsLen := len(expects)
	var expectStr string
	if expectsLen > 1 {
		for i, expect := range expects {
			if i < expectsLen-1 {
				expectStr += expect + ", "
			} else {
				expectStr += "or " + expect
			}
		}
	} else {
		expectStr = expects[0]
	}

	return fmt.Errorf("unexpected end of token stream, expecting %s", expectStr)
}

var ErrNoMoreToken = errors.New("no more token")
var ErrAlreadyAtBeginning = errors.New("already at beginning")
