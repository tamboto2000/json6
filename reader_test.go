package json6

import (
	"bytes"
	"testing"
)

func TestReader(t *testing.T) {
	input := `this string should have

	   5 new line`

	r := newReader(bytes.NewReader([]byte(input)), newPosition(1, 0))

	// read all the characters
	for {
		if _, _, err := r.ReadRune(); err != nil {
			break
		}
	}

	if r.p.ln != 5 {
		t.Errorf("unexpected %d lines, expecting %d", r.p.ln, 5)
	}
}
