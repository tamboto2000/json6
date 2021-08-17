package test

import (
	"fmt"
	"testing"

	"github.com/tamboto2000/json6"
)

func TestArrayValInterface(t *testing.T) {
	src := []byte("[1, 2,,,,, 3,,]")
	var v interface{}
	if err := json6.UnmarshalBytes(src, &v); err != nil {
		t.Error(err.Error())
		return
	}

	expectingLen := 3
	expectingIdxAndVal := map[int]int64{
		0: 1,
		1: 2,
		2: 3,
	}

	slice := v.([]interface{})
	if len(slice) != expectingLen {
		t.Errorf("Expecting len = %d, got %d instead", expectingLen, len(slice))
	}

	for i, v := range expectingIdxAndVal {
		if slice[i].(int64) != v {
			t.Errorf("Expecting value of index %d = %d, got %d instead", i, v, slice[i])
		}
	}
}

func TestArrayValSlice(t *testing.T) {
	src := []byte("[1, 2,,,,, 3,,]")
	var v []interface{}
	if err := json6.UnmarshalBytes(src, &v); err != nil {
		t.Error(err.Error())
		return
	}

	expectingLen := 3
	expectingIdxAndVal := map[int]int64{
		0: 1,
		1: 2,
		2: 3,
	}

	if len(v) != expectingLen {
		t.Errorf("Expecting len = %d, got %d instead", expectingLen, len(v))
	}

	for i, val := range expectingIdxAndVal {
		if v[i].(int64) != val {
			t.Errorf("Expecting value of index %d = %d, got %d instead", i, val, v[i])
		}
	}
}

func TestArrayValArray(t *testing.T) {
	src := []byte("[1, 2,,,,, 3,4,5]")
	var v [4]int
	if err := json6.UnmarshalBytes(src, &v); err != nil {
		t.Error(err.Error())
		return
	}

	expectingIdxAndVal := map[int]int{
		0: 1,
		1: 2,
		2: 3,
	}

	for i, val := range expectingIdxAndVal {
		if v[i] != val {
			t.Errorf("Expecting value of index %d = %d, got %d instead", i, val, v[i])
		}
	}
}

func TestArrayInvalid(t *testing.T) {
	src := []byte("[1, 2,,,,, 3,,]")
	var v int
	if err := json6.UnmarshalBytes(src, &v); err != nil {
		return
	}

	t.Error("Unmarshaler should return error")
}

func TestArrayNested(t *testing.T) {
	src := []byte("[[1], [1, 2], [1, 2, 3],,,[,,,],,,]")
	var v [][]int
	if err := json6.UnmarshalBytes(src, &v); err != nil {
		t.Error(err.Error())
		return
	}

	fmt.Printf("%#v", v)
}
