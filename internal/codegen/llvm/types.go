package llvm

import (
	"github.com/HicaroD/Telia/third/go-llvm"
)

type LLVMValue interface {
	Value() string
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
