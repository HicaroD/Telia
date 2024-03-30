package token

import "fmt"

type Position struct {
	Filename     string
	Line, Column int
}

func NewPosition(filename string, column, line int) Position {
	return Position{Filename: filename, Line: line, Column: column}
}

func (pos Position) String() string {
	return fmt.Sprintf("[%s:%d:%d]", pos.Filename, pos.Line, pos.Column)
}
