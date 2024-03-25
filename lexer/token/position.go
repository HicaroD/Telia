package token

import "fmt"

type Position struct {
	Filename string
	X, Y     int
}

func NewPosition(filename string, x, y int) Position {
	return Position{Filename: filename, X: x, Y: y}
}

func (pos Position) String() string {
	return fmt.Sprintf("[%s:%d:%d]", pos.Filename, pos.Y, pos.X)
}
