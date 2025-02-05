package main

import (
	"log"

	"github.com/HicaroD/Telia/backend/codegen/llvm"
	"github.com/HicaroD/Telia/diagnostics"
	"github.com/HicaroD/Telia/frontend/ast"
	"github.com/HicaroD/Telia/frontend/lexer"
	"github.com/HicaroD/Telia/frontend/parser"
	"github.com/HicaroD/Telia/middleend/sema"
)

func main() {
	args := cli()

	switch args.Command {
	case COMMAND_BUILD:
		var program *ast.Program
		var err error

		collector := diagnostics.New()

		if args.IsModuleBuild {
			program, err = buildModule(args, collector)
		} else {
			program, err = buildFile(args, collector)
		}

		// TODO(errors)
		if err != nil {
			log.Fatal(err)
		}

		sema := sema.New(collector)
		err = sema.Check(program)
		// TODO(errors)
		if err != nil {
			log.Fatal(err)
		}

		// TODO: define flag for setting the back-end
		// Currently I only have one type of back-end, but, in the future, I
		// could have more
		codegen := llvm.NewCG(args.Path)
		err = codegen.Generate(program)
		// TODO(errors)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func buildModule(cliResult CliResult, collector *diagnostics.Collector) (*ast.Program, error) {
	p := parser.New(collector)
	program, err := p.ParseModuleDir(cliResult.Path)
	if err != nil {
		return nil, err
	}

	return program, nil
}

func buildFile(cliResult CliResult, collector *diagnostics.Collector) (*ast.Program, error) {
	p := parser.New(collector)

	l, err := lexer.NewFromFilePath(cliResult.ParentDirName, cliResult.Path, collector)
	if err != nil {
		return nil, err
	}
	program, err := p.ParseFileAsProgram(l)
	if err != nil {
		return nil, err
	}

	return program, nil
}
