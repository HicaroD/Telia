package integration

import (
	"strings"
	"testing"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestSiblingLoopBodiesHaveIndependentDeclarations(t *testing.T) {
	for _, loop := range []string{
		"while count < 2",
		"for i := 0; i < 2; i = i + 1",
	} {
		t.Run(loop, func(t *testing.T) {
			path := loopControlSource(t, `package main
fn main() i32 {
count := 0
`+loop+` {
value := 1
count = count + value
}
count = 0
`+loop+` {
value := 1
count = count + value
}
if count != 2 {
return 1
}
return 0
}
`)
			runLoopControl(t, path)
		})
	}
}

func TestLoopDeclarationsDoNotEscapeTheirScope(t *testing.T) {
	for _, loop := range []string{
		"while false {\nvalue := 1\n}",
		"for i := 0; i < 0; i = i + 1 {\nvalue := 1\n}",
		"for value := 0; value < 0; value = value + 1 {\n}",
		"values i32[1] := .{1}\nfor value in values {\n}",
	} {
		t.Run(loop, func(t *testing.T) {
			path := loopControlSource(
				t,
				"package main\nfn main() i32 {\n"+loop+"\nreturn value\n}\n",
			)
			_, diags := compiler.CompileOnly(path)
			if len(diags.Diags) == 0 || !strings.Contains(diags.Diags[0].Message, "not declared") {
				t.Fatalf("want inaccessible loop-local diagnostic, got %v", diags.Diags)
			}
		})
	}
}

func TestDuplicateDeclarationWithinLoopBodyIsRejected(t *testing.T) {
	path := loopControlSource(t, `package main
fn main() i32 {
while false {
value := 1
value := 2
}
return 0
}
`)
	_, diags := compiler.CompileOnly(path)
	if len(diags.Diags) != 1 || !strings.Contains(diags.Diags[0].Message, "already declared") {
		t.Fatalf("want duplicate declaration diagnostic, got %v", diags.Diags)
	}
}

func TestLoopsResolveDeclarationsInEnclosingAnonymousBlock(t *testing.T) {
	path := loopControlSource(t, `package main
fn main() i32 {
{
limit := 2
count := 0
while count < limit {
count = count + 1
}
for i := 0; i < limit; i = i + 1 {
count = count + 1
}
if count != 4 {
return 1
}
}
return 0
}
`)
	runLoopControl(t, path)
}
