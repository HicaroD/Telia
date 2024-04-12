package values

import "tinygo.org/x/go-llvm"

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
