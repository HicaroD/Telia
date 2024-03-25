package parser

import (
	"github.com/HicaroD/telia-lang/lexer/token"
	"github.com/HicaroD/telia-lang/lexer/token/kind"
)

type cursor struct {
	offset int
	tokens []token.Token
}

func newCursor(tokens []token.Token) *cursor {
	return &cursor{offset: 0, tokens: tokens}
}

func (cursor *cursor) peek() *token.Token {
	if cursor.isOutOfBound() {
		return nil
	}
	return &cursor.tokens[cursor.offset]
}

func (cursor *cursor) next() *token.Token {
	if cursor.isOutOfBound() {
		return nil
	}
	token := cursor.tokens[cursor.offset]
	cursor.offset++
	return &token
}

func (cursor *cursor) skip() {
	cursor.next()
}

func (cursor *cursor) nextIs(expectedKind kind.TokenKind) bool {
	token := cursor.peek()
	if token == nil {
		return false
	}
	return token.Kind == expectedKind
}

func (cursor *cursor) isOutOfBound() bool {
	return cursor.offset >= len(cursor.tokens)
}
