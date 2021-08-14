package test

import (
	"testing"

	"github.com/tamboto2000/json6"
)

func TestHexa(t *testing.T) {
	src := []byte(`0x123abc`)
	var v int
	if err := json6.UnmarshalBytes(src, &v); err != nil {
		t.Error(err.Error())
	}

	expect := 1194684
	if v != expect {
		t.Error("Expecting", expect, "but got", v)
	}
}

func TestHexaWithUpperX(t *testing.T) {
	src := []byte(`0X123abc`)
	var v int
	if err := json6.UnmarshalBytes(src, &v); err != nil {
		t.Error(err.Error())
	}

	expect := 1194684
	if v != expect {
		t.Error("Expecting", expect, "but got", v)
	}
}
