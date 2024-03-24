package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/HicaroD/telia-lang/lexer"
)

func main() {
	args := os.Args[1:]
	filename := args[0]

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("unable to open file: %s due to error '%s'", filename, err)
	}

	reader := bufio.NewReader(file)

	lex := lexer.NewLexer(filename, reader)
	tokens := lex.Tokenize()
	for _, token := range tokens {
		switch token.Lexeme.(type) {
		case int:
		case string:
			fmt.Printf("%s '%s' %s\n", token.Kind, token.Lexeme, token.Position)
		default:
			fmt.Printf("%s %s\n", token.Kind, token.Position)
		}
	}
}
