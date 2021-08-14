package test

import (
	"testing"

	"github.com/tamboto2000/json6"
)

func TestOctal1(t *testing.T) {
	src := []byte(`0o123`)
	var v int
	if err := json6.UnmarshalBytes(src, &v); err != nil {
		t.Error(err.Error())
	}

	expecting := 0o123
	if v != expecting {
		t.Error("Expecting", expecting, "but got", v)
	}
}

func TestOctal2(t *testing.T) {
	src := []byte(`0O123`)
	var v int
	if err := json6.UnmarshalBytes(src, &v); err != nil {
		t.Error(err.Error())
	}

	expecting := 0O123
	if v != expecting {
		t.Error("Expecting", expecting, "but got", v)
	}
}

func TestOctal3(t *testing.T) {
	src := []byte(`0o1__2__3__`)
	var v int
	if err := json6.UnmarshalBytes(src, &v); err != nil {
		t.Error(err.Error())
	}

	expecting := 0o123
	if v != expecting {
		t.Error("Expecting", expecting, "but got", v)
	}
}
