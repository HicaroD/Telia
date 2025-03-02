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
)

type CliResult struct {
	Command   Command
	BuildType config.BuildType
	Loc       *ast.Loc
}

func cli() (CliResult, error) {
	result := CliResult{}

	args := os.Args[1:]
	if len(args) == 0 {
		return result, fmt.Errorf("TODO: show help")
	}

	command := args[0]

	switch command {
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

		result.Loc = ast.LocFromPath(fileOrDirPath)
	default:
		return result, fmt.Errorf("TODO: show help")
	}

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
	return result, nil
}
