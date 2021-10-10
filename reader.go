package json6

import "io"

type reader struct {
	p        *Position
	r        io.RuneScanner
	lastChar rune
}

func newReader(r io.RuneScanner, pos *Position) *reader {
	return &reader{p: pos, r: r}
}

func (r *reader) ReadRune() (rune, int, error) {
	char, size, err := r.r.ReadRune()
	if err != nil {
		return 0, 0, err
	}

	r.lastChar = char
	if char == '\n' {
		if r.lastChar != '\r' {
			r.p.addLn(1)
			r.p.setCol(0)
		}
	} else if char == '\r' {
		r.p.addLn(1)
		r.p.setCol(0)
	} else {
		r.p.addCol(1)
	}

	return char, size, nil
}

func (r *reader) UnreadRune() error {

	return nil
}
