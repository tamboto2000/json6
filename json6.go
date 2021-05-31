// Package json6 is a JSON6 implementation.
// Refer to here 
package json6

const (
	null = iota
	undefined
	boolean
	str
	integer
	float
	obj
	array
)

type object struct {
	kind uint
	boolean bool
	str string
	integer int64
	float float64
	obj map[string]*object
	array []*object
}
