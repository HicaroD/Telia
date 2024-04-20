package main

import (
	"bufio"
	"log"
	"os"

	"github.com/HicaroD/Telia/codegen"
	"github.com/HicaroD/Telia/lexer"
	"github.com/HicaroD/Telia/parser"
	"github.com/HicaroD/Telia/sema"
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
	defer file.Close()

	reader := bufio.NewReader(file)
	lex := lexer.New(filename, reader)
	tokens, err := lex.Tokenize()
	if err != nil {
		log.Fatal(err)
	}

	parser := parser.New(tokens)
	astNodes, err := parser.Parse()
	// TODO(errors)
	if err != nil {
		log.Fatal(err)
	}

	sema := sema.New()
	err = sema.Analyze(astNodes)
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
