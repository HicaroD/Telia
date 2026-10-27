// Package source owns immutable source snapshots and half-open byte spans.
package source

import (
	"sort"
	"strings"
)

type File struct {
	Path  string
	Text  string
	lines []int
}

func New(path string, text []byte) *File {
	f := &File{Path: path, Text: string(text), lines: []int{0}}
	for i, b := range text {
		if b == '\n' {
			f.lines = append(f.lines, i+1)
		}
	}
	return f
}

type Span struct {
	File       *File
	Start, End int
}

func (s Span) Valid() bool {
	return s.File != nil && s.Start >= 0 && s.End >= s.Start && s.End <= len(s.File.Text)
}

// Position returns one-based line and byte column, matching lexer positions.
func (f *File) Position(offset int) (int, int) {
	if offset < 0 {
		offset = 0
	}
	if offset > len(f.Text) {
		offset = len(f.Text)
	}
	i := sort.Search(len(f.lines), func(i int) bool { return f.lines[i] > offset }) - 1
	return i + 1, offset - f.lines[i] + 1
}

func (f *File) Line(line int) (string, int) {
	if line < 1 || line > len(f.lines) {
		return "", len(f.Text)
	}
	start := f.lines[line-1]
	end := len(f.Text)
	if line < len(f.lines) {
		end = f.lines[line] - 1
	}
	return strings.TrimSuffix(f.Text[start:end], "\r"), start
}
