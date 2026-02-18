# Telia Compiler Overview

Telia is a custom programming language compiler written in Go that targets LLVM for code generation. It features a hand-written lexer, recursive descent parser, semantic analyzer, and LLVM-based code generator to produce native executables.

## Project Structure

```
Telia/
├── cmd/compiler/           # Entry point and CLI
│   ├── main.go            # Main function
│   └── cli.go             # Command-line interface
├── internal/
│   ├── ast/               # Abstract Syntax Tree definitions
│   │   ├── ast.go         # Node kinds and base types
│   │   ├── types.go       # Type system definitions
│   │   └── scope.go       # Scope management
│   ├── lexer/             # Lexical analysis
│   │   ├── lexer.go       # Tokenizer implementation
│   │   └── token/         # Token definitions
│   │       └── kind.go    # Token kinds and keywords
│   ├── parser/            # Parsing (recursive descent)
│   │   └── parser.go      # Parser implementation
│   ├── sema/             # Semantic analysis
│   │   └── sema.go        # Type checking and inference
│   ├── codegen/           # Code generation
│   │   └── llvm/          # LLVM IR generation
│   │       └── llvm.go    # Codegen implementation
│   └── diagnostics/       # Error reporting
├── base/
│   ├── std/              # Standard library
│   │   ├── io/           # Input/output functions
│   │   ├── libc/         # C library bindings
│   │   └── math/         # Math functions
│   └── runtime/          # Runtime support
│       ├── mem.t         # Memory operations
│       └── deref.t       # Pointer dereference checks
├── examples/             # Example Telia programs
├── docs/                 # Documentation
└── Makefile             # Build commands
```

## Dependencies

- **Go 1.22+** - Compiler implementation language
- **LLVM 18** - Code generation backend
  - `opt` - LLVM optimizer
  - `clang-18` - C compiler for linking

## Build System

The Makefile provides the following commands:

| Command | Description |
|---------|-------------|
| `make build` | Build the Telia compiler binary |
| `make clean` | Remove build artifacts |
| `make test` | Run tests |
| `make fmt` | Format code |
| `make run-example` | Build and run example programs |

### Build Output

The compiler produces native executables through:
1. Generating LLVM IR (`.ll` file)
2. Running LLVM optimizer (`opt -O3`)
3. Linking with clang (`clang-18 -lm`)

## Language Features

Telia is a systems programming language with:

- **Functions**: `fn name(params) return_type { ... }`
- **Structs**: `struct Name { field type }`
- **Control Flow**: `if/elif/else`, `for`, `while`
- **Pointers**: `*Type` for pointers, `&expr` for address-of
- **FFI**: `extern` blocks for C interop
- **Defer**: `defer` statement for cleanup

### Type System

- **Basic Types**: `bool`, `i8` to `i128`, `u8` to `u128`, `f32`, `f64`, `string`, `cstring`, `rawptr`
- **Compound Types**: Pointers (`*T`), Structs, Tuples
- **Type Inference**: Variables can omit explicit types with `:=`

## Compiler Pipeline

```
Source (.t) → Lexer → Tokens → Parser → AST → Sema → Typed AST → Codegen → LLVM IR → Executable
```

Each stage transforms the program representation and performs specific analysis or generation tasks.
