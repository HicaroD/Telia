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

// NOTE: this variable is set during build
var DevMode string

var HELP_COMMAND string = `Telia - A simple, powerful, and flexible programming language for modern applications.
Telia offers robust features for building high-performance applications with simplicity and flexibility.

Usage:
  telia <command> [arguments]

Available Commands:
  build [path] [-release] [-debug]   Builds the program
      [path]        Path to the directory or file (defaults to current directory)
      -release      Build in release mode
      -debug        Build in debug mode

  env                               Show environment information

  help                              Show this help message

Examples:
  telia build                        Build the program in the current directory
  telia build path/to/project        Build the program in the specified directory
  telia build myfile.t -debug        Build the program in debug mode (or just omit the flag)
  telia build myfile.t -release      Build the program in release mode
  telia env                          Display environment details

For more information about Telia, visit: https://github.com/HicaroD/Telia
`

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
		config.ENVS.ShowAll()
		return
	case COMMAND_BUILD:
		collector := diagnostics.New()

		buildType := BUILD_TYPE_FILE
		if args.Loc.IsPackage {
			buildType = BUILD_TYPE_PACKAGE
		}
		program, runtime, err := buildAll(args.ArgLoc, args.Loc, collector, buildType)

		// TODO(errors)
		if err != nil {
			log.Fatal(err)
		}

		sema := sema.New(collector)
		err = sema.Check(program, runtime)
		// TODO(errors)
		if err != nil {
			log.Fatal(err)
		}

		// TODO: define flag for setting the back-end
		// Currently I only have one type of back-end, but, in the future, I
		// could have more

		// TODO: properly set directory
		codegen := llvm.NewCG(args.Loc, program, runtime)
		err = codegen.Generate(args.BuildOptType)
		// TODO(errors)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func SetupAll() error {
	config.SetDevMode(DevMode == "1")
	if config.DEV {
		fmt.Println("[DEV] initialized")
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

type BuildType int

const (
	BUILD_TYPE_FILE BuildType = iota
	BUILD_TYPE_PACKAGE
)

func (bt BuildType) String() string {
	switch bt {
	case BUILD_TYPE_FILE:
		return "file"
	case BUILD_TYPE_PACKAGE:
		return "package"
	}
	return "unknown"
}

func buildAll(
	argLoc string,
	loc *ast.Loc,
	collector *diagnostics.Collector,
	ty BuildType,
) (*ast.Program, *ast.Package, error) {
	switch ty {
	case BUILD_TYPE_PACKAGE:
		p := parser.New(collector)
		program, runtime, err := p.ParsePackageAsProgram(argLoc, loc)
		return program, runtime, err
	case BUILD_TYPE_FILE:
		p := parser.New(collector)
		program, runtime, err := p.ParseFileAsProgram(argLoc, loc, collector)
		return program, runtime, err
	default:
		return nil, nil, fmt.Errorf("unknown build type: %s", ty)
	}
}
