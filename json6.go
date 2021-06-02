// Package json6 is a JSON6 implementation.
// Refer to https://github.com/d3x0r/JSON6 for more info
package json6

const (
	null = iota
	undefined
	boolean
	backTickstr
	doubleQuoteStr
	singleQuoteStr
	integer
	float
	obj
	array
	comment
)

type object struct {
	kind    uint
	boolean bool
	str     string
	integer int64
	float   float64
	obj     map[string]*object
	array   []*object
	rns     []rune
}
