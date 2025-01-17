package main

import (
	"fmt"

	"github.com/HicaroD/Telia/backend/codegen/llvm"
	"github.com/HicaroD/Telia/diagnostics"
	"github.com/HicaroD/Telia/frontend/parser"
	"github.com/HicaroD/Telia/middleend/sema"
)

func main() {
	args := cli()

	switch args.Command {
	case COMMAND_BUILD:
		if args.IsModuleBuild {
			err := buildModule(args)
			if err != nil {
				return
			}
		}
	}
}

func buildModule(cliResult CliResult) error {
	collector := diagnostics.New()

	p := parser.New(collector)
	program, err := p.ParseModuleDir(cliResult.Path)
	if err != nil {
		return err
	}

	fmt.Println(program)
	sema := sema.New(collector)

	err = sema.Check(program)
	if err != nil {
		return err
	}

	// TODO: define flag for setting the back-end
	// Currently I only have one type of back-end, but in the future
	// I could have more
	codegen := llvm.NewCG(cliResult.Path)
	err = codegen.Generate(program)
	return err
}
