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

func new(filename string, reader *bufio.Reader) *cursor {
	return &cursor{reader: reader, position: token.NewPosition(filename, 1, 1)}
}

func (cursor *cursor) peek() (rune, bool) {
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

func (cursor *cursor) skip() {
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

func (cursor *cursor) skipWhitespace() {
	cursor.readWhile(func(character rune) bool { return unicode.IsSpace(character) })
}

func (cursor *cursor) readWhile(isValid func(rune) bool) string {
	var content strings.Builder

	for {
		character, ok := cursor.peek()
		if !ok {
			break
		}
		if isValid(character) {
			content.WriteRune(character)
			cursor.skip()
		} else {
			break
		}
	}

	return content.String()
}

func (cursor cursor) Position() token.Position {
	return cursor.position
}
