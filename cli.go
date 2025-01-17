package main

import (
	"log"
	"os"
)

type Command int

const (
	COMMAND_BUILD Command = iota
)

type CliResult struct {
	Command Command

	IsModuleBuild bool   // true if 'Command' is build and 'Path' is directory
	Path          string // path to directory
}

func cli() CliResult {
	result := CliResult{}

	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatal("TODO: show help - list of commands")
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

		result.Path = fileOrDir
		result.IsModuleBuild = info.Mode().IsDir()
	default:
		log.Fatal("TODO: show help - list of commands")
	}

	return result
}
