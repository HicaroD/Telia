package codegen

import "tinygo.org/x/go-llvm"

type variable struct {
	ty  llvm.Type
	ptr llvm.Value
}

