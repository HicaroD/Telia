package main

import (
	"log"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/codegen/llvm"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/parser"
	"github.com/HicaroD/Telia/internal/sema"
)

func main() {
	args, err := cli()
	if err != nil {
		log.Fatal(err)
	}

	switch args.Command {
	case COMMAND_BUILD:
		var program *ast.Program
		var err error

		collector := diagnostics.New()

		if args.Loc.IsPackage {
			program, err = buildPackage(args.Loc, collector)
		} else {
			program, err = buildFile(args.Loc, collector)
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

		// TODO: properly set directory
		codegen := llvm.NewCG(args.Loc, program)
		err = codegen.Generate(args.BuildType)
		// TODO(errors)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func buildPackage(loc *ast.Loc, collector *diagnostics.Collector) (*ast.Program, error) {
	p := parser.New(collector)
	program, err := p.ParsePackageAsProgram(loc)
	return program, err
}

func buildFile(loc *ast.Loc, collector *diagnostics.Collector) (*ast.Program, error) {
	p := parser.New(collector)
	program, err := p.ParseFileAsProgram(loc, collector)
	return program, err
}
