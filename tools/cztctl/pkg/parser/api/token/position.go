package token

import "fmt"

// IllegalPosition is a position that is not valid.
var IllegalPosition = Position{}

// Position represents a rune position in the source code.
type Position struct {
	Filename string
	Line     int
	Column   int
}

// String returns a string representation of the position.
func (p Position) String() string {
	if len(p.Filename) == 0 {
		return fmt.Sprint(p.Line, ":", p.Column)
	}
	return fmt.Sprint(p.Filename, " ", p.Line, ":", p.Column)
}
