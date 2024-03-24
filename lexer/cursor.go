package lexer

import (
	"bufio"
)

type Cursor struct {
	reader *bufio.Reader
	x      int
	y      int
}

func newCursor(reader *bufio.Reader) *Cursor {
	return &Cursor{reader: reader, x: 0, y: 1}
}

func (cursor *Cursor) Next() (rune, error) {
	character, _, err := cursor.reader.ReadRune()
	if err != nil {
		return 0, err
	}

	if character == '\n' {
		cursor.x = 1
		cursor.y++
	} else {
		cursor.x++
	}

	return character, nil
}

func (cursor *Cursor) Peek() (rune, error) {
	var err error

	character, _, err := cursor.reader.ReadRune()
	if err != nil {
		return 0, err
	}

	err = cursor.reader.UnreadRune()
	if err != nil {
		return 0, err
	}

	return character, nil
}

func (cursor *Cursor) Position() (int, int) {
	return cursor.x, cursor.y
}
