# LLVM Code Generation

The code generator transforms the typed AST into LLVM Intermediate Representation (IR), then uses LLVM tools to produce a native executable.

## Overview

The code generator:
1. Creates an LLVM module
2. Emits declarations (functions, structs)
3. Emits function bodies (instructions)
4. Runs LLVM optimizer
5. Links with clang to produce executable

## LLVM Type Mapping

Telia types are mapped to LLVM types:

| Telia Type | LLVM Type |
|------------|-----------|
| `bool` | `i1` |
| `i8`, `u8` | `i8` |
| `i16`, `u16` | `i16` |
| `i32`, `u32`, `int`, `uint` | `i32` |
| `i64`, `u64` | `i64` |
| `i128`, `u128` | `i128` |
| `f32` | `float` |
| `f64` | `double` |
| `string`, `cstring` | `i8*` |
| `rawptr` | `i8*` |
| `*T` | `T*` |
| `struct { ... }` | `{ ... }` |
| `(T1, T2, ...)` | `{ T1, T2, ... }` |

The type mapping function translates Telia expression types to their corresponding LLVM representations.

## Code Generation Functions

### Program Generation

1. Generates runtime package
2. Generates main program package
3. Produces executable

### Package Generation

- Processes package files recursively
- Emits declarations first
- Emits bodies after declarations

### Declaration Emission

- Function signatures
- External declarations
- Struct types

### Function Emission

1. Creates entry basic block
2. Allocates parameter storage
3. Emits block statements

## Expression Code Generation

### Expression Emission

Each expression generates:
- Type of the expression (LLVM type)
- LLVM value (constant or instruction)
- Whether value contains float (for optimization decisions)

### Pointer-to-Value Conversion

When an expression is used as a value:
- For address-of expressions: don't load (we want the pointer)
- For literals: don't load (already a value)
- For function calls: don't load (return value)
- Otherwise: load from pointer to get the value

### Load/Store Pattern

Local variables are stored in alloca slots (stack-allocated):
- Allocate space for the variable
- Store the initial value into the allocated space
- Load from the alloca when reading the variable

## Control Flow Generation

### Conditionals

Creates basic blocks:
- `.if` - then branch
- `.else` - else branch  
- `.end` - after conditional

### For Loops

Creates basic blocks:
- `.forprep` - preparation
- `.forinit` - initialization
- `.forbody` - loop body
- `.forupdate` - update expression
- `.forend` - after loop

### While Loops

Creates basic blocks:
- `.whileinit` - condition check
- `.whilebody` - loop body
- `.whileend` - after loop

## Runtime Calls

The codegen emits calls to runtime functions for:
- Nil pointer dereference checks before pointer operations
- Other runtime support as needed

## Build Pipeline

### IR Generation

```
Source → AST → Sema → LLVM IR (.ll file)
```

### Optimization

```bash
opt -O3 input.ll -o optimized.ll
```

### Linking

```bash
clang-18 -o output optimized.ll -lm
```

Flags:
- `-Wl,-s` - Strip debug symbols (release build)
- `-O0` - No optimization (debug build)
- `-O3` - Maximum optimization (release build)

## Variable Backend Representation

Variables store both:
- Allocated type (LLVM type)
- Alloca pointer (memory location for load/store operations)

## Adding New Constructs

When adding new AST node kinds:

1. **Expression Emission**:
   - Add case for new expression kind
   - Return type, value, and float flag

2. **Statement Emission**:
   - Add case for new statement kind

3. **Type Mapping**:
   - Add mapping for new type if needed

4. **Operator Emission**:
   - Add emission logic for new operators
