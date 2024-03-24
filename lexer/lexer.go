package lexer

import (
	"bufio"
	"fmt"
	"io"

	"github.com/HicaroD/telia-lang/lexer/token"
	"github.com/HicaroD/telia-lang/lexer/token/kind"
)

// TODO: define everything related to the lexer, such as cursor
type Lexer struct {
	cursor *Cursor
}

func NewLexer(reader *bufio.Reader) *Lexer {
	return &Lexer{cursor: newCursor(reader)}
}

func (lex *Lexer) Tokenize() []token.Token {
	var tokens []token.Token
	// TODO: build tokenizer loop for reading all tokens
	for {
		character, err := lex.cursor.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		x, y := lex.cursor.Position()
		fmt.Printf("%c [%d:%d]\n", character, x, y)
		// tokens = append(tokens, token)
	}
	return tokens
}

func getToken() token.Token {
	return token.NewToken(kind.Id)
}
