// Package native resolves package-owned native requirements without network I/O.
package native

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Target struct {
	Libraries    []string `json:"libraries"`
	RuntimeFiles []string `json:"runtime_files"`
	IncludeDirs  []string `json:"include_dirs"`
	Sources      []string `json:"sources"`
	CompileFlags []string `json:"compile_flags"`
	LinkFlags    []string `json:"link_flags"`
}

type Manifest struct {
	Version string            `json:"version"`
	Targets map[string]Target `json:"targets"`
}

type Requirements struct {
	CompileArgs  []string
	LinkArgs     []string
	RuntimeFiles []string
}

func Resolve(manifestPath, target string) (Requirements, error) {
	var result Requirements
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return result, fmt.Errorf("native package: %w", err)
	}
	var manifest Manifest
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&manifest); err != nil {
		return result, fmt.Errorf("native manifest %s: %w", manifestPath, err)
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return result, fmt.Errorf("native manifest %s: unexpected trailing data", manifestPath)
	}
	if manifest.Version == "" {
		return result, fmt.Errorf("native manifest %s requires a pinned version", manifestPath)
	}
	entry, ok := manifest.Targets[target]
	if !ok {
		return result, fmt.Errorf(
			"native package %s does not support target %s",
			manifestPath,
			target,
		)
	}
	dir := filepath.Dir(manifestPath)
	for _, input := range []struct {
		paths     []string
		directory bool
	}{
		{entry.IncludeDirs, true}, {entry.Sources, false},
		{entry.Libraries, false}, {entry.RuntimeFiles, false},
	} {
		for _, path := range input.paths {
			if path == "" || filepath.IsAbs(path) || filepath.VolumeName(path) != "" ||
				strings.Contains(path, "\\") ||
				!filepath.IsLocal(path) {
				return result, fmt.Errorf("native path %q must be package-relative", path)
			}
			fullPath := filepath.Join(dir, path)
			info, err := os.Stat(fullPath)
			if err != nil {
				return result, fmt.Errorf(
					"missing native artifact for %s: %s; run the package setup script: %w",
					target,
					fullPath,
					err,
				)
			}
			if (input.directory && !info.IsDir()) ||
				(!input.directory && !info.Mode().IsRegular()) {
				return result, fmt.Errorf(
					"invalid native artifact for %s: %s; run the package setup script",
					target,
					fullPath,
				)
			}
		}
	}
	for _, include := range entry.IncludeDirs {
		result.CompileArgs = append(result.CompileArgs, "-I", filepath.Join(dir, include))
	}
	result.CompileArgs = append(result.CompileArgs, entry.CompileFlags...)
	for _, source := range entry.Sources {
		result.CompileArgs = append(result.CompileArgs, filepath.Join(dir, source))
	}
	for _, library := range entry.Libraries {
		result.LinkArgs = append(result.LinkArgs, filepath.Join(dir, library))
	}
	result.LinkArgs = append(result.LinkArgs, entry.LinkFlags...)
	for _, file := range entry.RuntimeFiles {
		result.RuntimeFiles = append(result.RuntimeFiles, filepath.Join(dir, file))
	}
	return result, nil
}
