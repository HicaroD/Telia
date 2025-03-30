package main

import (
	"fmt"
	"log"

	"github.com/HicaroD/Telia/config"
	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/codegen/llvm"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/parser"
	"github.com/HicaroD/Telia/internal/sema"
)

// NOTE: this variable is set to 1 during build
var DevMode string

func main() {
	err := SetupAll()
	if err != nil {
		log.Fatal(err)
	}

	args, err := cli()
	if err != nil {
		log.Fatal(err)
	}

	switch args.Command {
	case COMMAND_HELP:
		fmt.Print(HELP_COMMAND)
		return
	case COMMAND_ENV:
		for k, v := range config.ENVS {
			fmt.Printf("%s='%s'\n", k, v)
		}
		return
	case COMMAND_BUILD:
		var program *ast.Program
		var err error

		collector := diagnostics.New()

		if args.Loc.IsPackage {
			program, err = buildPackage(args.ArgLoc, args.Loc, collector)
		} else {
			program, err = buildFile(args.ArgLoc, args.Loc, collector)
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

func SetupAll() error {
	config.SetDevMode(DevMode == "1")
	if config.DEV {
		fmt.Println("[DEV MODE] initialized")
	}

	err := config.SetupConfigDir()
	if err != nil {
		return err
	}

	err = config.SetupEnvFile()
	if err != nil {
		return err
	}

	return nil
}

func buildPackage(
	argLoc string,
	loc *ast.Loc,
	collector *diagnostics.Collector,
) (*ast.Program, error) {
	p := parser.New(collector)
	program, err := p.ParsePackageAsProgram(argLoc, loc)
	return program, err
}

func buildFile(
	argLoc string,
	loc *ast.Loc,
	collector *diagnostics.Collector,
) (*ast.Program, error) {
	p := parser.New(collector)
	program, err := p.ParseFileAsProgram(argLoc, loc, collector)
	return program, err
}
