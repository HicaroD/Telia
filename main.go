package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/HicaroD/Telia/backend/codegen"
	"github.com/HicaroD/Telia/diagnostics"
	"github.com/HicaroD/Telia/frontend/lexer"
	"github.com/HicaroD/Telia/frontend/parser"
	"github.com/HicaroD/Telia/middleend/sema"
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

	diagCollector := diagnostics.New()
	reader := bufio.NewReader(file)

	lex := lexer.New(filename, reader, diagCollector)
	// tokens, err := lex.Tokenize()
	// if err != nil {
	// 	fmt.Println("Errors found during compilation ", err)
	// 	os.Exit(1)
	// }

	parser := parser.New(lex, diagCollector)
	// astNodes, err := parser.Parse()
	// if err != nil {
	// 	fmt.Println("Errors found during compilation: ", err)
	// 	os.Exit(1)
	// }

	sema := sema.New(parser, diagCollector)
	err = sema.Analyze(astNodes)
	if err != nil {
		fmt.Println("Errors found during compilation: ", err)
		os.Exit(1)
	}

	codegen := codegen.New(filename, astNodes)
	err = codegen.Generate()
	// TODO(errors)
	if err != nil {
		log.Fatal(err)
	}
}
