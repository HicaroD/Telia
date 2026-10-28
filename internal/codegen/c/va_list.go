package c

import (
	"fmt"

	"github.com/HicaroD/Telia/internal/ast"
)

func vaCall(node *ast.Node) *ast.FnCall {
	for node.Kind == ast.KIND_NAMESPACE_ACCESS {
		node = node.Node.(*ast.NamespaceAccess).Right
	}
	call, _ := node.Node.(*ast.FnCall)
	return call
}

func (c *CCodegen) emitVaOwnerInitialization(variable *ast.VarIdStmt, expr *ast.Node) {
	name := variable.BackendType.(*CVariable).Name
	call := vaCall(expr)
	origin := fmt.Sprintf("va_list owner left scope without va_end (%s:%d, %s)", variable.Name.Pos.Filename, variable.Name.Pos.Line, variable.Name.Name())
	fmt.Fprintf(c.buf, "_TeliaVaList %s __attribute__((cleanup(_telia_va_scope_end))) = {0};\n%s.origin = %s;\n", name, name, emitStringLiteral([]byte(origin)))
	if call.Builtin == "va_start" {
		last := call.VaLast.BackendType.(*CVariable).Name
		fmt.Fprintf(c.buf, "va_start(%s.args, %s);\n%s.state = 1;\n", name, last, name)
		return
	}
	var source string
	if borrow := vaCall(call.Args[0]); borrow != nil && borrow.Builtin == "va_borrow" {
		original := borrow.VaOwner.BackendType.(*CVariable).Name
		fmt.Fprintf(c.buf, "if (%s.state != 1) _panic(%s);\n", original, vaPanicMessage(borrow, "va_list is ended or consumed"))
		source = original + ".args"
	} else {
		source = c.emitExpr(call.Args[0])
	}
	fmt.Fprintf(c.buf, "va_copy(%s.args, %s);\n%s.state = 1;\n", name, source, name)
}

func (c *CCodegen) emitVaOperation(call *ast.FnCall) string {
	owner := call.VaOwner.BackendType.(*CVariable).Name
	if call.Builtin == "va_arg" {
		return fmt.Sprintf("({ if (%s.state != 1) _panic(%s); va_arg(%s.args, %s); })", owner, vaPanicMessage(call, "va_list is ended or consumed"), owner, emitCType(call.BuiltinType))
	}
	if call.Builtin == "va_end" {
		return fmt.Sprintf("({ if (%s.state == 0) _panic(%s); va_end(%s.args); %s.state = 0; })", owner, vaPanicMessage(call, "va_list has already ended"), owner, owner)
	}
	return fmt.Sprintf("({ if (%s.state != 1) _panic(%s); %s.state = 2; %s.args; })", owner, vaPanicMessage(call, "va_list is ended or consumed"), owner, owner)
}

func vaPanicMessage(call *ast.FnCall, message string) string {
	pos := call.Name.Pos
	return emitStringLiteral([]byte(fmt.Sprintf("%s (%s:%d:%d)", message, pos.Filename, pos.Line, pos.Column)))
}
