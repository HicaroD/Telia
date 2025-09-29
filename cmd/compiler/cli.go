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
		result.BuildOptType = config.BUILD_OPT_DEBUG

		// TODO: check for unknown flags, the user might mispell the command
		for _, arg := range args[1:] {
			switch arg {
			case "-release":
				if releaseBuildSet {
					return result, fmt.Errorf("duplicate -release flag")
				}
				releaseBuildSet = true
				result.BuildOptType = config.BUILD_OPT_RELEASE
			case "-debug":
				if debugBuildSet {
					return result, fmt.Errorf("duplicate -debug flag")
				}
				debugBuildSet = true
				result.BuildOptType = config.BUILD_OPT_DEBUG
			}
		}
		// TODO(errors)
		if releaseBuildSet && debugBuildSet {
			return result, fmt.Errorf("choose either -release or -build, not both")
		}
	case "help":
		fallthrough
	default:
		result.Command = COMMAND_HELP
		return result, nil
	}
	return result, nil
}
