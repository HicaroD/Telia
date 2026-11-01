package c

import _ "embed"

//go:embed runtime.h
var runtimeCore string

//go:embed runtime_os.h
var runtimeOS string

//go:embed runtime_process.h
var runtimeProcess string

//go:embed runtime_hash.h
var runtimeHash string

//go:embed runtime_seed.h
var runtimeSeed string

var runtimeSource = runtimeCore + "\n" + runtimeOS + "\n" + runtimeProcess + "\n" + runtimeHash + "\n" + runtimeSeed
