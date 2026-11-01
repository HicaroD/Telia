package integration

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/HicaroD/Telia/tests/compiler"
)

// Complete documentation examples pair a telia fence with exact stdout. Other
// fences are fragments or existing examples covered by their own integration tests.
func TestDocumentationExamples(t *testing.T) {
	root := filepath.Join("..", "..", "docs")
	pattern := regexp.MustCompile(
		"(?s)```telia\\n(.*?)\\n```\\n\\nExpected output:\\n\\n```text\\n(.*?)\\n```",
	)
	count := 0
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		text := strings.ReplaceAll(string(data), "\r\n", "\n")
		examples := pattern.FindAllStringSubmatch(text, -1)
		if len(examples) != strings.Count(text, "```telia\n") {
			t.Errorf("%s: every telia fence must have an expected-output block", path)
		}
		for index, example := range examples {
			count++
			name := strings.TrimPrefix(path, root+string(filepath.Separator))
			t.Run(fmt.Sprintf("%s/%d", name, index+1), func(t *testing.T) {
				source := filepath.Join(t.TempDir(), "main.t")
				if err := os.WriteFile(source, []byte(example[1]+"\n"), 0600); err != nil {
					t.Fatal(err)
				}
				for _, opt := range []string{"-O0", "-O3"} {
					t.Run(opt, func(t *testing.T) {
						binary, diags := compiler.CompileOnlyWithCFlags(source, opt)
						if diags.HasErrors() {
							t.Fatalf("compile documentation example: %v", diags.Diags)
						}
						defer removeTestBinary(t, binary)
						ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
						defer cancel()
						cmd := exec.CommandContext(ctx, binary)
						var stderr strings.Builder
						cmd.Stderr = &stderr
						output, err := cmd.Output()
						if err != nil {
							t.Fatalf("run: %v; stderr: %s", err, stderr.String())
						}
						if got, want := string(output), example[2]+"\n"; got != want {
							t.Fatalf("want stdout %q, got %q", want, got)
						}
					})
				}
			})
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("no complete documentation examples found")
	}
}
