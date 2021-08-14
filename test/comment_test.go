package test

import (
	"testing"

	"github.com/tamboto2000/json6"
)

func TestInlineComment(t *testing.T) {
	src := []byte("// inline comment")
	var v string
	if err := json6.UnmarshalBytes(src, &v); err != nil {
		t.Error(err.Error())
	}
}

func TestMultiLineComment(t *testing.T) {
	src := []byte(`
	/* 
	multi-line comment
	 */
	`)
	var v string
	if err := json6.UnmarshalBytes(src, &v); err != nil {
		t.Error(err.Error())
	}
}
