package main

import (
	"bufio"
	"log"
	"os"

	"github.com/HicaroD/telia-lang/codegen"
	"github.com/HicaroD/telia-lang/lexer"
	"github.com/HicaroD/telia-lang/parser"
	"github.com/HicaroD/telia-lang/sema"
)

func main() {
	args := os.Args[1:]
	// TODO(errors)
	if len(args) == 0 {
		log.Fatal("error: no input files")
	}
	filename := args[0]

	file, err := os.Open(filename)
	// TODO(errors)
	if err != nil {
		log.Fatalf("unable to open file: %s due to error '%s'", filename, err)
	}

	reader := bufio.NewReader(file)
	lex := lexer.New(filename, reader)
	tokens := lex.Tokenize()

	parser := parser.New(tokens)
	astNodes, err := parser.Parse()
	// TODO(errors)
	if err != nil {
		log.Fatal(err)
	}

	sema := sema.New(astNodes)
	err = sema.Analyze()
	// TODO(errors)
	if err != nil {
		log.Fatal(err)
	}

	codegen := codegen.New(astNodes)
	err = codegen.Generate()
	// TODO(errors)
	if err != nil {
		log.Fatal(err)
	}
}
