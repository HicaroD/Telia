package values

import (
	"tinygo.org/x/go-llvm"
)

type Function struct {
	LLVMValue
	Fn     llvm.Value
	Ty     llvm.Type
	Locals map[string]*LLVMValue
}

func NewFunctionValue(fn llvm.Value, ty llvm.Type) Function {
	return Function{Fn: fn, Ty: ty, Locals: map[string]*LLVMValue{}}
}

func (function Function) Value() string {
	return "Function"
}
