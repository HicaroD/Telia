package codegen

import "tinygo.org/x/go-llvm"

type cache struct {
	globals   map[string]*llvm.Value
	functions map[string]*function
}

func newCache() *cache {
	return &cache{
		globals:   map[string]*llvm.Value{},
		functions: map[string]*function{},
	}
}

func (cache *cache) InsertFunction(name string, fn *llvm.Value, ty *llvm.Type) *function {
	// TODO(errors)
	if fn, ok := cache.functions[name]; ok {
		return fn
	}
	function := function{
		fn:     fn,
		ty:     ty,
		locals: map[string]*variable{},
	}
	cache.functions[name] = &function
	return &function
}

func (cache *cache) GetFunction(name string) *function {
	if fn, ok := cache.functions[name]; ok {
		return fn
	}
	return nil
}

func (cache *cache) InsertGlobal(name string, value *llvm.Value) {
	// TODO(errors)
	if _, ok := cache.globals[name]; ok {
		return
	}
	cache.globals[name] = value
}

func (cache *cache) GetGlobal(name string) *llvm.Value {
	if global, ok := cache.globals[name]; ok {
		return global
	}
	return nil
}
