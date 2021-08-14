package test

import (
	"testing"

	"github.com/tamboto2000/json6"
)

func TestDecodeNaN(t *testing.T) {
	src := []byte(`NaN`)
	var v float64
	if err := json6.UnmarshalBytes(src, &v); err != nil {
		t.Error(err.Error())
	}
}
