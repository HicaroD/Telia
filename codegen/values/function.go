package values

import (
	"tinygo.org/x/go-llvm"
)

type Function struct {
	LLVMValue
	Fn     llvm.Value
	Ty     llvm.Type
	Block  *llvm.BasicBlock
	Locals map[string]*LLVMValue
}

func NewFunctionValue(fn llvm.Value, ty llvm.Type, block *llvm.BasicBlock) *Function {
	return &Function{Fn: fn, Ty: ty, Block: block, Locals: map[string]*LLVMValue{}}
}

func (function Function) Value() string {
	return "Function"
}

