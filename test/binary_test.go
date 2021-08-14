package test

import (
	"testing"

	"github.com/tamboto2000/json6"
)

func TestBinary1(t *testing.T) {
	src := []byte(`0b10101010`)
	var v int
	if err := json6.UnmarshalBytes(src, &v); err != nil {
		t.Error(err.Error())
	}

	expecting := 170
	if v != expecting {
		t.Error("Expecting", expecting, "but got", v)
	}
}

func TestBinary2(t *testing.T) {
	src := []byte(`0B10101010`)
	var v int
	if err := json6.UnmarshalBytes(src, &v); err != nil {
		t.Error(err.Error())
	}

	expecting := 170
	if v != expecting {
		t.Error("Expecting", expecting, "but got", v)
	}
}

func TestBinary3(t *testing.T) {
	src := []byte(`0b1__0__101010__`)
	var v int
	if err := json6.UnmarshalBytes(src, &v); err != nil {
		t.Error(err.Error())
	}

	expecting := 170
	if v != expecting {
		t.Error("Expecting", expecting, "but got", v)
	}
}

func TestBinary4(t *testing.T) {
	src := []byte(`-0b1__0__101010__`)
	var v int
	if err := json6.UnmarshalBytes(src, &v); err != nil {
		t.Error(err.Error())
	}

	expecting := -170
	if v != expecting {
		t.Error("Expecting", expecting, "but got", v)
	}
}
