package test

import (
	"testing"

	"github.com/tamboto2000/json6"
)

func TestOnlyDigits1(t *testing.T) {
	src := []byte(`123`)
	var v int
	if err := json6.UnmarshalBytes(src, &v); err != nil {
		t.Error(err.Error())
	}

	expecting := 123
	if v != expecting {
		t.Error("Expecting", expecting, "but got", v)
	}
}

func TestOnlyDigits2(t *testing.T) {
	src := []byte(`1__2__3__`)
	var v int
	if err := json6.UnmarshalBytes(src, &v); err != nil {
		t.Error(err.Error())
	}

	expecting := 123
	if v != expecting {
		t.Error("Expecting", expecting, "but got", v)
	}
}

func TestOnlyDigits3(t *testing.T) {
	src := []byte(`00123`)
	var v int
	if err := json6.UnmarshalBytes(src, &v); err != nil {
		t.Error(err.Error())
	}

	expecting := 123
	if v != expecting {
		t.Error("Expecting", expecting, "but got", v)
	}
}

func TestOnlyDigits4(t *testing.T) {
	src := []byte(`123`)
	var v interface{}
	if err := json6.UnmarshalBytes(src, &v); err != nil {
		t.Error(err.Error())
	}

	expecting := 123
	if v.(int64) != int64(expecting) {
		t.Error("Expecting", expecting, "but got", v)
	}
}
