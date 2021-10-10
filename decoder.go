package json6

import "reflect"

// ValueType define value type of a JSON6 value
type ValueType uint

// value types
const (
	ValueNull ValueType = iota
	ValueBoolean
	ValueString
	ValueInteger
	ValueFloat
	ValueObject
	ValueArray
)

// Value contains decoded value from token or sequence of tokens (like array and objects)
type Value struct {
	t      ValueType
	refVal reflect.Value
	valMap map[string]Value // if t == ValueObject
	vals   []Value          // if t == ValueArray
}

// Decoder decode tokens into JSON6 value
type Decoder struct {
}
