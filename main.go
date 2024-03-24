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
	for i := range tokens {
		fmt.Printf("%s %s\n", tokens[i].Kind, tokens[i].Position)
	}
}
