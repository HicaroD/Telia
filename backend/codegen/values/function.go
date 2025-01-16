package values

import (
	"tinygo.org/x/go-llvm"
)

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
