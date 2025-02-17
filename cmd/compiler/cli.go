package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/HicaroD/Telia/internal/config"
)

type Command int

const (
	COMMAND_BUILD Command = iota
)

type CliResult struct {
	Command Command

	IsModuleBuild bool   // true if 'Command' is build and 'Path' is directory
	ParentDirName string // name of parent dir
	Path          string // path to directory / file (treated as module)

	BuildType config.BuildType
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

		fileOrDir := "."
		if len(args) >= 2 {
			fileOrDir = args[1]
		}

		info, err := os.Stat(fileOrDir)
		if err != nil {
			log.Fatalf("os.Stat error: %s\n", err)
		}

		path, err := filepath.Abs(fileOrDir)
		if err != nil {
			log.Fatal(err)
		}

		result.Path = path
		result.IsModuleBuild = info.Mode().IsDir()
		if result.IsModuleBuild {
			result.ParentDirName = filepath.Base(path)
		} else {
			result.ParentDirName = filepath.Base(filepath.Dir(path))
		}
	default:
		return result, fmt.Errorf("TODO: show help")
	}

	releaseBuild, debugBuild := false, false
	for _, arg := range args[1:] {
		switch arg {
		case "-release":
			releaseBuild = true
			result.BuildType = config.RELEASE
		case "-debug":
			debugBuild = true
			result.BuildType = config.DEBUG
		}
	}

	if releaseBuild && debugBuild {
		return result, fmt.Errorf("choose either -release or -build, not both")
	}
	if !releaseBuild && !debugBuild {
		result.BuildType = config.DEBUG
	}

	return result, nil
}
