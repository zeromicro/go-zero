package token

import "fmt"

var IllegalPosition = Position{}

type Position struct {
	Filename string
	Line     int
	Column   int
}

func (p Position) String() string {
	if len(p.Filename) == 0 {
		return fmt.Sprint(p.Line, ":", p.Column)
	}
	return fmt.Sprint(p.Filename, " ", p.Line, ":", p.Column)
}
