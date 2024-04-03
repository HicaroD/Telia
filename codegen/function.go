package codegen

import (
	"log"

	"tinygo.org/x/go-llvm"
)

type function struct {
	fn     *llvm.Value
	ty     *llvm.Type
	locals map[string]*variable
}

func (function *function) InsertLocal(name string, ty llvm.Type, ptr llvm.Value) error {
	// TODO(errors)
	if _, ok := function.locals[name]; ok {
		log.Fatalf("this should not happen at code generation level - variable redeclaration")
	}
	variable := variable{
		ty:  ty,
		ptr: ptr,
	}
	function.locals[name] = &variable
	return nil
}

func (function *function) GetLocal(name string) *variable {
	if local, ok := function.locals[name]; ok {
		return local
	}
	return nil
}
