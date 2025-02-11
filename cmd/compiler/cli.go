package main

import (
	"log"
	"os"
	"path/filepath"
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

		path, err := filepath.Abs(fileOrDir)
		if err != nil {
			log.Fatal(err)
		}

		result.Path = path
		result.IsModuleBuild = info.Mode().IsDir()
		result.ParentDirName = filepath.Base(filepath.Dir(path))
	default:
		log.Fatal("TODO: show help - list of commands")
	}

	return result
}
