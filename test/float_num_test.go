package test

import (
	"testing"

	"github.com/tamboto2000/json6"
)

func TestFloatNum1(t *testing.T) {
	src := []byte(`0.123`)
	var v float64
	if err := json6.UnmarshalBytes(src, &v); err != nil {
		t.Error(err.Error())
	}

	expect := 0.123
	if v != expect {
		t.Error("Expecting", expect, "but got", v)
	}
}

func TestFloatNum2(t *testing.T) {
	src := []byte(`-0.123`)
	var v float64
	if err := json6.UnmarshalBytes(src, &v); err != nil {
		t.Error(err.Error())
	}

	expect := -0.123
	if v != expect {
		t.Error("Expecting", expect, "but got", v)
	}
}

func TestFloatNum3(t *testing.T) {
	src := []byte(`-.123`)
	var v float64
	if err := json6.UnmarshalBytes(src, &v); err != nil {
		t.Error(err.Error())
	}

	expect := -0.123
	if v != expect {
		t.Error("Expecting", expect, "but got", v)
	}
}

func TestFloatNum4(t *testing.T) {
	src := []byte(`-.__1__2__3__`)
	var v float64
	if err := json6.UnmarshalBytes(src, &v); err != nil {
		t.Error(err.Error())
	}

	expect := -0.123
	if v != expect {
		t.Error("Expecting", expect, "but got", v)
	}
}

func TestFloatNum5(t *testing.T) {
	// this should be producing error	
	src := []byte(`-._______`)
	var v float64
	if err := json6.UnmarshalBytes(src, &v); err != nil {
		return
	}
}

func TestFloatNum6(t *testing.T) {
	// this should be producing error	
	src := []byte(`-.e_123_____`)
	var v float64
	if err := json6.UnmarshalBytes(src, &v); err != nil {
		return
	}
}

func TestFloatNumNonZero(t *testing.T) {
	src := []byte(`1.123`)
	var v float64
	if err := json6.UnmarshalBytes(src, &v); err != nil {
		t.Error(err.Error())
	}

	expect := 1.123
	if v != expect {
		t.Error("Expecting", expect, "but got", v)
	}
}
