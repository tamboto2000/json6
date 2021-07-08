package json6

import (
	"bytes"
	"io"
	"strings"
	"text/scanner"
)

func scanFromBytes(byts []byte) *scanner.Scanner {
	r := bytes.NewReader(byts)
	s := new(scanner.Scanner)
	s.Init(r)

	return s
}

func scanFromString(str string) *scanner.Scanner {
	r := strings.NewReader(str)
	s := new(scanner.Scanner)
	s.Init(r)

	return s
}

func scanFromReader(r io.Reader) *scanner.Scanner {
	s := new(scanner.Scanner)
	s.Init(r)

	return s
}
