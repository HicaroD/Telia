package token

import "fmt"

type Position struct {
	Filename string
	Column, Line     int
}

func NewPosition(filename string, x, y int) Position {
	return Position{Filename: filename, Column: x, Line: y}
}

func (pos Position) String() string {
	return fmt.Sprintf("[%s:%d:%d]", pos.Filename, pos.Line, pos.Column)
}
