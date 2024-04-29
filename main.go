package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/HicaroD/Telia/codegen"
	"github.com/HicaroD/Telia/collector"
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

	diagCollector := collector.New()
	reader := bufio.NewReader(file)

	lex := lexer.New(filename, reader, diagCollector)
	tokens, err := lex.Tokenize()
	if err != nil {
		fmt.Println("Errors found during compilation ", err)
		os.Exit(1)
	}

	parser := parser.New(tokens, diagCollector)
	astNodes, err := parser.Parse()
	if err != nil {
		fmt.Println("Errors found during compilation: ", err)
		os.Exit(1)
	}

	sema := sema.New(diagCollector)
	err = sema.Analyze(astNodes)
	if err != nil {
		fmt.Println("Errors found during compilation: ", err)
		os.Exit(1)
	}

	codegen := codegen.New(astNodes)
	err = codegen.Generate()
	// TODO(errors)
	if err != nil {
		log.Fatal(err)
	}
}
