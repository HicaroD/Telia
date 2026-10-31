package sema

import (
	"fmt"
	"strings"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/lexer/token"
	"github.com/HicaroD/Telia/internal/source"
)

type matchHintCase struct {
	subject  *ast.IdExpr
	patterns []string
	keys     []string
	guard    string
}

func guardedMatchCase(n *ast.Node) *matchHintCase {
	b, ok := n.Node.(*ast.BinExpr)
	if !ok {
		return nil
	}
	if b.Op == token.AND || b.Op == token.OR {
		left, right := guardedMatchCase(b.Left), guardedMatchCase(b.Right)
		if left == nil || right == nil || left.subject.N != right.subject.N {
			return nil
		}
		return &matchHintCase{subject: left.subject, guard: hintSource(n)}
	}
	if b.Op <= token.CMP_OP_START || b.Op >= token.CMP_OP_END || b.OperandType == nil ||
		!b.OperandType.IsInteger() {
		return nil
	}
	id, ok := b.Left.Node.(*ast.IdExpr)
	if !ok || !stableMatchSubject(id) {
		return nil
	}
	// Constant evaluation alone is insufficient: only literal syntax is safe
	// here, not arbitrary expressions or foreign constants.
	literal := b.Right
	if unary, ok := literal.Node.(*ast.UnaryExpr); ok && unary.Op == token.MINUS {
		literal = unary.Value
	}
	value, ok := literal.Node.(*ast.LiteralExpr)
	if !ok || !value.Type.IsInteger() {
		return nil
	}
	return &matchHintCase{subject: id, guard: hintSource(n)}
}

func stableMatchSubject(id *ast.IdExpr) bool {
	return id.N != nil && (id.N.Kind == ast.KIND_PARAM || id.N.Kind == ast.KIND_VAR_ID_STMT)
}

// Only suggest matching a stable local value, never an expression whose
// evaluation count could change (calls, field reads, or pointer dereferences).
func equalityMatchCase(n *ast.Node) *matchHintCase {
	b, ok := n.Node.(*ast.BinExpr)
	if !ok {
		return nil
	}
	if b.Op == token.OR {
		left, right := equalityMatchCase(b.Left), equalityMatchCase(b.Right)
		if left == nil || right == nil || left.subject.N != right.subject.N {
			return nil
		}
		left.patterns = append(left.patterns, right.patterns...)
		left.keys = append(left.keys, right.keys...)
		return left
	}
	if b.Op != token.EQUAL_EQUAL {
		return nil
	}
	id, ok := b.Left.Node.(*ast.IdExpr)
	if !ok || !stableMatchSubject(id) {
		return nil
	}
	var key, text string
	switch value := b.Right.Node.(type) {
	case *ast.LiteralExpr:
		if !value.Type.IsInteger() {
			return nil
		}
		constant, err := evaluateConstant(b.Right)
		if err != nil {
			return nil
		}
		key = constant.ExactString()
		text = string(value.Value)
	case *ast.NamespaceAccess:
		if value.EnumVariant == nil || len(value.EnumVariant.Fields) != 0 {
			return nil
		}
		key = fmt.Sprintf("%p", value.EnumVariant)
		text = strings.Trim(hintSource(b.Right), "() \t\r\n")
	default:
		return nil
	}
	if text == "" {
		return nil
	}
	return &matchHintCase{subject: id, patterns: []string{text}, keys: []string{key}}
}

func hintSource(n *ast.Node) string {
	if n == nil || !n.Span.Valid() {
		return ""
	}
	return n.Span.File.Text[n.Span.Start:n.Span.End]
}

func (s *sema) suggestMatch(cond *ast.CondStmt) {
	span := cond.IfStmt.Expr.Span
	if !span.Valid() || s.matchHints[span] {
		return
	}
	branches := append([]*ast.IfElifCond{cond.IfStmt}, cond.ElifStmts...)
	// A short chain is already clear. Only offer actionable rewrites, not
	// general style advice for every conditional or range check.
	if len(branches) < 3 || cond.ElseStmt == nil {
		return
	}
	cases := make([]*matchHintCase, 0, len(branches))
	seen := map[string]bool{}
	equalities := 0
	for _, branch := range branches {
		c := equalityMatchCase(branch.Expr)
		if c == nil {
			c = guardedMatchCase(branch.Expr)
		}
		if c == nil || (len(cases) > 0 && c.subject.N != cases[0].subject.N) {
			return
		}
		if len(c.patterns) > 0 {
			equalities++
		}
		for _, key := range c.keys {
			if seen[key] {
				return // A literal rewrite would introduce a duplicate pattern.
			}
			seen[key] = true
		}
		cases = append(cases, c)
	}
	if equalities < 2 {
		return
	}
	name := cases[0].subject.Name.Name()
	var rewrite strings.Builder
	fmt.Fprintf(&rewrite, "return match %s {\n", name)
	for i, branch := range branches {
		value := hintReturnValue(branch.Block)
		if value == "" {
			return
		}
		pattern := strings.Join(cases[i].patterns, " | ")
		if cases[i].guard != "" {
			pattern = "_ if " + cases[i].guard
		}
		fmt.Fprintf(&rewrite, "    %s => %s\n", pattern, value)
	}
	fallback := hintReturnValue(cond.ElseStmt.Block)
	if fallback == "" {
		return
	}
	fmt.Fprintf(&rewrite, "    _ => %s\n}", fallback)
	if s.matchHints == nil {
		s.matchHints = map[source.Span]bool{}
	}
	s.matchHints[span] = true
	s.collector.ReportAndSave(diagnostics.Diag{
		Severity: diagnostics.Hint,
		Title:    "Consider using match",
		Message: fmt.Sprintf(
			"These branches test `%s`. A match expression can make the cases easier to read.",
			name,
		),
		Span: cond.IfStmt.Expr.Span,
		Help: "Suggested rewrite:\n\n" + rewrite.String(),
	})
}

func hintReturnValue(block *ast.BlockStmt) string {
	if len(block.Statements) != 1 {
		return ""
	}
	ret, ok := block.Statements[0].Node.(*ast.ReturnStmt)
	if !ok {
		return ""
	}
	return hintSource(ret.Value)
}
