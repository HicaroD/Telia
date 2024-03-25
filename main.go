package main

import (
	"bufio"
	"log"
	"os"

	// "github.com/HicaroD/telia-lang/ast"
	"github.com/HicaroD/telia-lang/lexer"
	"github.com/HicaroD/telia-lang/parser"
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

	parser := parser.NewParser(tokens)
	// astNodes, err := parser.Parse()
	_, err = parser.Parse()
	if err != nil {
		// TODO(errors)
		log.Fatal(err)
	}

	// for i := range astNodes {
	// 	switch astNodes[i].(type) {
	// 	case ast.FunctionDecl:
	// 		fnDecl := astNodes[i].(ast.FunctionDecl)
	// 	}
	// }
}
