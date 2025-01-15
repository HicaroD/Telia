package token

import "fmt"

// TODO: since I changed the architecture to support modules and files,
// I probably don't need to store the filename on every
// position, saving memory!
type Pos struct {
	Filename     string
	Line, Column int
}

func NewPosition(filename string, column, line int) Pos {
	return Pos{Filename: filename, Line: line, Column: column}
}

func (pos *Pos) Move(character byte) {
	if character == '\n' {
		pos.Column = 1
		pos.Line++
	} else {
		pos.Column++
	}
}

func (pos *Pos) SetPosition(newPos Pos) {
	// TODO: is there a better way to replace this object
	pos.Filename = newPos.Filename
	pos.Line = newPos.Line
	pos.Column = newPos.Column
}

func (pos Pos) String() string {
	return fmt.Sprintf("[%s:%d:%d]", pos.Filename, pos.Line, pos.Column)
}
