package main

import (
	"fmt"
	"os"

	"github.com/HicaroD/Telia/config"
	"github.com/HicaroD/Telia/internal/ast"
	ccodegen "github.com/HicaroD/Telia/internal/codegen/c"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/parser"
	"github.com/HicaroD/Telia/internal/sema"
	"github.com/HicaroD/Telia/internal/source"
	"golang.org/x/term"
)

// NOTE: this variable is set during build
var DevMode string

var HELP_COMMAND string = `Telia - A simple, powerful, and flexible programming language for modern applications.
Telia offers robust features for building high-performance applications with simplicity and flexibility.

Usage:
  telia <command> [arguments]

Available Commands:
  build [path] [-release] [-debug] [-o <path>]   Builds the program
      [path]        Path to the directory or file (defaults to current directory)
      -release      Build in release mode
      -debug        Build in debug mode
      -o <path>     Output executable path (default: source name in current directory)
      --color=auto|always|never  Diagnostic colors (default: auto; respects NO_COLOR)

  env                               Show environment information

  help                              Show this help message

Examples:
  telia build                        Build the program in the current directory
  telia build path/to/project        Build the program in the specified directory
  telia build myfile.t -debug        Build the program in debug mode (or just omit the flag)
  telia build myfile.t -release      Build the program in release mode
  telia build myfile.t -o myapp      Build and output to 'myapp'
  telia build myfile.t -o bin/myapp  Build and output to 'bin/myapp'
  telia env                          Display environment details

For more information about Telia, visit: https://github.com/HicaroD/Telia
`

func main() {
	os.Exit(run())
}

func run() int {
	err := SetupAll()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	args, err := cli()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	switch args.Command {
	case COMMAND_HELP:
		fmt.Print(HELP_COMMAND)
		return 0
	case COMMAND_ENV:
		config.ENVS.ShowAll()
		return 0
	case COMMAND_BUILD:
		collector := diagnostics.New()
		_, noColor := os.LookupEnv("NO_COLOR")
		renderer := diagnostics.Renderer{Color: diagnostics.ColorEnabled(args.Color, term.IsTerminal(int(os.Stderr.Fd())), noColor)}
		finish := func(err error) int {
			if !collector.HasErrors() {
				collector.ReportError(err, source.Span{})
			}
			if renderErr := renderer.Render(os.Stderr, collector); renderErr != nil {
				fmt.Fprintln(os.Stderr, renderErr)
				return 1
			}
			if err != nil || collector.HasErrors() {
				return 1
			}
			return 0
		}

		buildType := BUILD_TYPE_FILE
		if args.Loc.IsPackage {
			buildType = BUILD_TYPE_PACKAGE
		}
		program, err := buildAll(args.ArgLoc, args.Loc, collector, buildType)

		if err != nil || collector.HasErrors() {
			return finish(err)
		}

		sema := sema.New(collector)
		err = sema.Check(program)
		if err != nil || collector.HasErrors() {
			return finish(err)
		}

		codegen := ccodegen.NewCG(args.Loc, program)
		if args.OutputPath != "" {
			codegen.SetOutputPath(args.OutputPath)
		}
		err = codegen.Generate(args.BuildOptType)
		return finish(err)
	}
	return 0
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
) (*ast.Program, error) {
	switch ty {
	case BUILD_TYPE_PACKAGE:
		p := parser.New(collector)
		program, err := p.ParsePackageAsProgram(argLoc, loc)
		return program, err
	case BUILD_TYPE_FILE:
		p := parser.New(collector)
		program, err := p.ParseFileAsProgram(argLoc, loc, collector)
		return program, err
	default:
		return nil, fmt.Errorf("unknown build type: %s", ty)
	}
}
