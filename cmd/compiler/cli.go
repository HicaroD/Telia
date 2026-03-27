package main

import (
	"fmt"
	"log"
	"os"

	"github.com/HicaroD/Telia/config"
	"github.com/HicaroD/Telia/internal/ast"
)

type Command int

const (
	COMMAND_BUILD Command = iota
	COMMAND_HELP
	COMMAND_ENV
)

type CliResult struct {
	Command      Command
	BuildOptType config.BuildOptimizationType
	ArgLoc       string
	Loc          *ast.Loc
	OutputPath   string
}

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
	case "build":
		result.Command = COMMAND_BUILD

		releaseBuildSet, debugBuildSet := false, false
		result.BuildOptType = config.BUILD_OPT_DEBUG

		var sourcePath string

		for i := 1; i < len(args); i++ {
			arg := args[i]
			switch {
			case arg == "-release":
				if releaseBuildSet {
					return result, fmt.Errorf("duplicate -release flag")
				}
				releaseBuildSet = true
				result.BuildOptType = config.BUILD_OPT_RELEASE
			case arg == "-debug":
				if debugBuildSet {
					return result, fmt.Errorf("duplicate -debug flag")
				}
				debugBuildSet = true
				result.BuildOptType = config.BUILD_OPT_DEBUG
			case arg == "-o":
				i++
				if i >= len(args) {
					return result, fmt.Errorf("-o requires a path argument")
				}
				result.OutputPath = args[i]
			default:
				if sourcePath == "" {
					sourcePath = arg
				} else {
					return result, fmt.Errorf("unexpected argument: %s", arg)
				}
			}
		}

		if releaseBuildSet && debugBuildSet {
			return result, fmt.Errorf("choose either -release or -debug, not both")
		}

		if sourcePath == "" {
			sourcePath = "."
		}

		_, err := os.Stat(sourcePath)
		if err != nil {
			log.Fatalf("No such file or directory: %s\n", sourcePath)
		}

		loc, err := ast.LocFromPath(sourcePath)
		if err != nil {
			return result, err
		}
		result.ArgLoc = sourcePath
		result.Loc = loc
	case "help":
		fallthrough
	default:
		result.Command = COMMAND_HELP
		return result, nil
	}
	return result, nil
}
