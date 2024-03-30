package lexer

import (
	"bufio"
	"strings"
	"unicode"

	"github.com/HicaroD/telia-lang/lexer/token"
)

type cursor struct {
	reader   *bufio.Reader
	position token.Position
}

func newCursor(filename string, reader *bufio.Reader) *cursor {
	return &cursor{reader: reader, position: token.NewPosition(filename, 0, 0)}
}

func (cursor *cursor) Peek() (rune, bool) {
	var err error

	character, _, err := cursor.reader.ReadRune()
	// TODO(errors)
	if err != nil {
		return 0, false
	}

	err = cursor.reader.UnreadRune()
	// TODO(errors)
	if err != nil {
		return 0, false
	}

	return character, true
}

func (cursor *cursor) Skip() {
	// QUESTION: should I ignore the error here?
	cursor.next()
}

func (cursor *cursor) next() (rune, bool) {
	character, _, err := cursor.reader.ReadRune()
	// TODO(errors)
	if err != nil {
		return 0, false
	}

	if character == '\n' {
		cursor.position.Column = 1
		cursor.position.Line++
	} else {
		cursor.position.Column++
	}

	return character, false
}

func (cursor *cursor) SkipWhitespace() {
	cursor.ReadWhile(func(character rune) bool { return unicode.IsSpace(character) })
}

func (cursor *cursor) ReadWhile(isValid func(rune) bool) string {
	var content strings.Builder

	for {
		character, ok := cursor.Peek()
		if !ok {
			break
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

func (cursor cursor) Position() token.Position {
	return cursor.position
}
