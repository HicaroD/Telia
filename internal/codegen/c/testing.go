package c

import (
	"fmt"
	"strings"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/lexer/token"
)

const testingRuntime = `
static uint64_t _telia_test_failures = 0;
static _String _telia_test_name = {(const unsigned char *)"assertion", 9};
static void _telia_test_result(const char *label) {
    printf("%s ", label);
    fwrite(_telia_test_name.data, 1, _telia_test_name.len, stdout);
    fputc('\n', stdout);
    fflush(stdout);
}
static void _telia_testing_fail(void) {
    fputs("assertion failed\n", stderr);
    exit(1);
}
static void _telia_test_integer(unsigned __int128 value, bool negative) {
    char digits[40];
    int n = 0;
    do { digits[n++] = '0' + value % 10; value /= 10; } while (value);
    if (negative) fputc('-', stderr);
    while (n) fputc(digits[--n], stderr);
}
`

func (c *CCodegen) emitTestingAssertion(call *ast.FnCall) string {
	values := make([]string, len(call.Args))
	for i, arg := range call.Args {
		value := c.emitExpr(arg)
		values[i] = c.nextTmp()
		fmt.Fprintf(c.buf, "%s %s = %s;\n", emitCType(call.BuiltinType), values[i], value)
	}
	condition := values[0]
	if len(values) == 2 {
		if call.BuiltinType.Equals(ast.NewBasicType(token.STRING_TYPE)) {
			condition = fmt.Sprintf("_telia_string_equal(%s, %s)", values[0], values[1])
		} else {
			condition = values[0] + " == " + values[1]
		}
		if call.Builtin == "testing_assert_ne" {
			condition = "!(" + condition + ")"
		}
	}
	fmt.Fprintf(c.buf, "if (!(%s)) {\n", condition)
	c.emitTestFailureHeader(call.Name)
	if len(values) == 2 {
		label := "expected"
		if strings.HasSuffix(call.Builtin, "assert_ne") {
			label = "unexpected"
		}
		fmt.Fprintf(c.buf, "fputs(\"  %s: \", stderr);\n", label)
		c.emitTestValue(values[1], call.BuiltinType)
		c.buf.WriteString("fputs(\"\\n  actual: \", stderr);\n")
		c.emitTestValue(values[0], call.BuiltinType)
		c.buf.WriteString("fputc('\\n', stderr);\n")
	}
	c.emitTestFailureExit()
	c.buf.WriteString("}\n")
	return "((void)0)"
}

func (c *CCodegen) emitTestFailureHeader(site *token.Token) {
	fmt.Fprintf(
		c.buf,
		"fprintf(stderr, \"%%s:%%d:%%d: assertion failed in \", %s, %d, %d);\n",
		emitCBytes([]byte(site.Pos.Filename)),
		site.Pos.Line,
		site.Pos.Column,
	)
	c.buf.WriteString(
		"fwrite(_telia_test_name.data, 1, _telia_test_name.len, stderr); fputc('\\n', stderr);\n",
	)
}

func (c *CCodegen) emitTestFailureExit() {
	if c.program.TestMode {
		c.buf.WriteString("_telia_test_failures++;\n")
		c.emitTestUnwind()
	} else {
		c.buf.WriteString("_telia_testing_fail();\n")
	}
}

func (c *CCodegen) emitTestValue(value string, ty *ast.ExprType) {
	kind := ty.T.(*ast.BasicType).Kind
	switch {
	case kind == token.STRING_TYPE:
		fmt.Fprintf(
			c.buf,
			"fputc('\"', stderr); if (%s.len) fwrite(%s.data, 1, %s.len, stderr); fputc('\"', stderr);\n",
			value,
			value,
			value,
		)
	case kind == token.BOOL_TYPE:
		fmt.Fprintf(c.buf, "fputs(%s ? \"true\" : \"false\", stderr);\n", value)
	case kind.IsInteger():
		if kind >= token.UINT_TYPE && kind <= token.U128_TYPE {
			fmt.Fprintf(c.buf, "_telia_test_integer((unsigned __int128)%s, false);\n", value)
		} else {
			fmt.Fprintf(c.buf, "_telia_test_integer(%s < 0 ? 0 - (unsigned __int128)%s : (unsigned __int128)%s, %s < 0);\n", value, value, value, value)
		}
	default:
		fmt.Fprintf(c.buf, "fprintf(stderr, \"%%.17g\", (double)%s);\n", value)
	}
}

func (c *CCodegen) emitTestRunner() {
	c.buf.WriteString("\nint main(void) {\nint failed = 0;\nuint64_t before;\n")
	pass, fail := "PASS", "FAIL"
	if c.testColor {
		pass, fail = "\033[32mPASS\033[0m", "\033[31mFAIL\033[0m"
	}
	count := 0
	for _, file := range c.program.Root.Files {
		for _, node := range file.Body {
			if node.Kind != ast.KIND_FN_DECL {
				continue
			}
			fn := node.Node.(*ast.FnDecl)
			if fn.TestDescription == nil {
				continue
			}
			fmt.Fprintf(
				c.buf,
				"_telia_test_name = %s;\nbefore = _telia_test_failures;\n%s();\nif (_telia_test_failures == before) {\n_telia_test_result(%s);\n} else {\nfailed++;\n_telia_test_result(%s);\n}\n",
				emitStringLiteral([]byte(*fn.TestDescription)),
				fn.CName,
				emitCBytes([]byte(pass)),
				emitCBytes([]byte(fail)),
			)
			count++
		}
	}
	if count == 0 {
		c.buf.WriteString("puts(\"no tests found\");\n")
	}
	fmt.Fprintf(
		c.buf,
		"printf(\"%%d passed, %%d failed\\n\", %d - failed, failed);\nreturn failed ? 1 : 0;\n}\n",
		count,
	)
}

// A failure returns normally through Telia frames, running their registered
// cleanup. While emitting a deferred action, a new failure skips only that
// action's remainder; the surrounding cleanup sequence still runs.
func (c *CCodegen) emitTestUnwind() {
	previous := c.testUnwinding
	c.testUnwinding = true
	defer func() { c.testUnwinding = previous }()
	depth := 0
	if c.testCleanupLabel != "" {
		depth = c.testCleanupDepth
	}
	for i := len(c.deferScopes) - 1; i >= depth; i-- {
		c.flushDefers(i, "")
	}
	if c.testCleanupLabel != "" {
		fmt.Fprintf(c.buf, "goto %s;\n", c.testCleanupLabel)
	} else if c.currentRetTy == nil || c.currentRetTy.Equals(ast.NewBasicType(token.VOID_TYPE)) {
		c.buf.WriteString("return;\n")
	} else {
		fmt.Fprintf(c.buf, "return (%s){0};\n", emitCType(c.currentRetTy))
	}
}

func (c *CCodegen) emitCheckedTestCall(call *ast.FnCall) string {
	return c.vaExpression(func() string {
		before := c.nextTmp()
		fmt.Fprintf(c.buf, "uint64_t %s = _telia_test_failures;\n", before)
		value := c.emitFnCallExpr(call)
		ty := call.ResultType
		if call.Decl != nil {
			ty = call.Decl.RetType
		} else if call.Proto != nil {
			ty = call.Proto.RetType
		}
		result := "((void)0)"
		if ty == nil || ty.Equals(ast.NewBasicType(token.VOID_TYPE)) {
			fmt.Fprintf(c.buf, "%s;\n", value)
		} else {
			result = c.nextTmp()
			fmt.Fprintf(c.buf, "%s %s = %s;\n", emitCType(ty), result, value)
		}
		fmt.Fprintf(c.buf, "if (_telia_test_failures != %s) {\n", before)
		c.emitTestUnwind()
		c.buf.WriteString("}\n")
		return result
	})
}
