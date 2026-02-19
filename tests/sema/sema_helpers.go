package sema_test

import (
	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/lexer"
	"github.com/HicaroD/Telia/internal/parser"
	"github.com/HicaroD/Telia/internal/sema"
)

func parseAndCheck(src string) *diagnostics.Collector {
	collector := diagnostics.New()

	loc := &ast.Loc{Name: "test.t"}
	lex := lexer.New(loc, []byte(src), collector)
	p := parser.NewWithLex(lex, collector)

	file, err := p.ParseFileForTest(loc)
	if err != nil {
		diag := diagnostics.Diag{Message: "parse error: " + err.Error()}
		collector.ReportAndSave(diag)
		return collector
	}

	pkg := p.GetPkg()
	pkg.Loc = loc
	pkg.Files = []*ast.File{file}

	prog := &ast.Program{Root: pkg}

	s := sema.New(collector)
	err = s.Check(prog, nil)
	if err != nil {
		diag := diagnostics.Diag{Message: err.Error()}
		collector.ReportAndSave(diag)
	}

	return collector
}

func parseNextDecl(p *parser.Parser, file *ast.File) (*ast.Node, bool, error) {
	return p.NextForTest(file)
}

func containsDiag(diags []diagnostics.Diag, substr string) bool {
	for _, d := range diags {
		if len(d.Message) > 0 && containsString(d.Message, substr) {
			return true
		}
	}
	return false
}

func containsString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
