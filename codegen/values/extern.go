package values

import (
	"github.com/HicaroD/Telia/scope"
)

type Extern struct {
	LLVMValue
	Scope      *scope.Scope[LLVMValue]
}

func NewExtern(scope *scope.Scope[LLVMValue]) *Extern {
	return &Extern{Scope: scope}
}
