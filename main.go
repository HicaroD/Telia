package main

import (
	"fmt"
	"log"
	"os"

	"github.com/HicaroD/Telia/backend/codegen/llvm"
	"github.com/HicaroD/Telia/diagnostics"
	"github.com/HicaroD/Telia/frontend/lexer"
	"github.com/HicaroD/Telia/frontend/parser"
	"github.com/HicaroD/Telia/middleend/sema"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatal("error: no input files")
	}
	filename := args[0]

	// TODO: find more efficient ways to read a file into the memory
	// since calling os.ReadFile(...) reads the entire file into the memory
	// and this will cause problems in the future
	collector := diagnostics.New()

	file, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("unable to open file: %s due to error '%s'", filename, err)
	}

	lex := lexer.New(filename, file, collector)
	parser := parser.New(lex, collector)
	program, err := parser.Parse()
	if err != nil {
		fmt.Println("Errors found during compilation: ", err)
		os.Exit(1)
	}

	sema := sema.New(collector)
	err = sema.Check(program)
	if err != nil {
		fmt.Println("Errors found during compilation: ", err)
		os.Exit(1)
	}

	// TODO: define flag for setting the back-end
	// Currently I only have one type of back-end, but in the future
	// I could have more
	codegen := llvm.NewCG(filename)
	err = codegen.Generate(program)
	// TODO(errors)
	if err != nil {
		log.Fatal(err)
	}
}
