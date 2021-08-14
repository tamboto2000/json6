package test

import (
	"testing"

	"github.com/tamboto2000/json6"
)

func TestExponent(t *testing.T) {
	src := []byte(`0e123`)
	var v float64
	if err := json6.UnmarshalBytes(src, &v); err != nil {
		t.Error(err.Error())
	}

	expecting := 0e123
	if v != expecting {
		t.Error("Expecting", expecting, "but got", v)
	}
}

func TestDecimalExponent(t *testing.T) {
	src := []byte(`0.1e123`)
	var v float64
	if err := json6.UnmarshalBytes(src, &v); err != nil {
		t.Error(err.Error())
	}

	expecting := 0.1e123

	if v != expecting {
		t.Error("Expecting", 0, "but got", v)
	}
}

func TestExponentNonZero(t *testing.T) {
	src := []byte(`1e123`)
	var v float64
	if err := json6.UnmarshalBytes(src, &v); err != nil {
		t.Error(err.Error())
	}

	expect := 1e123
	if v != expect {
		t.Error("Expecting", expect, "but got", v)
	}
}
