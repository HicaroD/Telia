package main

import (
	"fmt"
	"log"
	"os"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/config"
)

type Command int

const (
	COMMAND_BUILD Command = iota
	COMMAND_HELP
	COMMAND_ENV
)

type CliResult struct {
	Command   Command
	BuildType config.BuildType
	ArgLoc    string
	Loc       *ast.Loc
}

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

func cli() (CliResult, error) {
	result := CliResult{}

	args := os.Args[1:]
	if len(args) == 0 {
		result.Command = COMMAND_HELP
		return result, nil
	}

	command := args[0]
	switch command {
	case "env":
		result.Command = COMMAND_ENV
	case "help":
		result.Command = COMMAND_HELP
	case "build":
		result.Command = COMMAND_BUILD

		fileOrDirPath := "."
		if len(args) >= 2 {
			fileOrDirPath = args[1]
		}

		_, err := os.Stat(fileOrDirPath)
		// TODO(errors)
		if err != nil {
			log.Fatalf("No such file or directory: %s\n", fileOrDirPath)
		}

		loc, err := ast.LocFromPath(fileOrDirPath)
		if err != nil {
			return result, err
		}
		result.ArgLoc = fileOrDirPath
		result.Loc = loc

		releaseBuildSet, debugBuildSet := false, false
		result.BuildType = config.DEBUG

		// TODO: check for unknown flags, the user might mispell the command
		for _, arg := range args[1:] {
			switch arg {
			case "-release":
				releaseBuildSet = true
				result.BuildType = config.RELEASE
			case "-debug":
				debugBuildSet = true
				result.BuildType = config.DEBUG
			}
		}
		// TODO(errors)
		if releaseBuildSet && debugBuildSet {
			return result, fmt.Errorf("choose either -release or -build, not both")
		}
	default:
		return result, fmt.Errorf("TODO: show help")
	}
	return result, nil
}
