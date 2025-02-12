package llvm

import (
	"tinygo.org/x/go-llvm"
)

type LLVMValue interface {
	Value() string
}

type Function struct {
	LLVMValue
	Fn llvm.Value
	Ty llvm.Type
}

func NewFunctionValue(fn llvm.Value, ty llvm.Type, block *llvm.BasicBlock) *Function {
	return &Function{Fn: fn, Ty: ty}
}

func (function Function) Value() string {
	return "Function"
}

type Variable struct {
	LLVMValue
	Ty  llvm.Type
	Ptr llvm.Value
}

func NewVariableValue(ty llvm.Type, ptr llvm.Value) *Variable {
	return &Variable{Ty: ty, Ptr: ptr}
}

func (variable Variable) Value() string {
	return "Variable"
}
