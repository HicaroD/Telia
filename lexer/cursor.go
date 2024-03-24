package lexer

import (
	"bufio"
	"io"
	"strings"
	"unicode"

	"github.com/HicaroD/telia-lang/lexer/token"
)

type Cursor struct {
	reader   *bufio.Reader
	position token.Position
}

func newCursor(filename string, reader *bufio.Reader) *Cursor {
	return &Cursor{reader: reader, position: token.NewPosition(filename, 0, 0)}
}

func (cursor *Cursor) Next() (rune, error) {
	character, _, err := cursor.reader.ReadRune()
	if err != nil {
		return 0, err
	}

	if character == '\n' {
		cursor.position.X = 1
		cursor.position.Y++
	} else {
		cursor.position.X++
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

func (cursor *Cursor) Skip() {
	// QUESTION: should I ignore the error here?
	cursor.Next()
}

func (cursor *Cursor) SkipWhitespace() {
	cursor.ReadWhile(func(character rune) bool { return unicode.IsSpace(character) })
}

func (cursor *Cursor) ReadWhile(isValid func(rune) bool) string {
	var content strings.Builder

	for {
		character, err := cursor.Peek()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		if isValid(character) {
			content.WriteRune(character)
			cursor.Skip()
		} else {
			break
		}
	}

	return content.String()
}

func (cursor Cursor) Position() token.Position {
	return cursor.position
}
