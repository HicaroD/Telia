# Compilation Pipeline

The Telia compiler processes source code through a series of stages, each transforming the program representation and performing specific analysis or generation tasks.

## Stage Overview

```
Source (.t) → Lexer → Parser → Semantic Analyzer → Codegen → Executable
```

### Stage Responsibilities

| Stage | Input | Output | Responsibilities |
|-------|-------|--------|-----------------|
| Lexer | Source text | Tokens | Character → token conversion, position tracking |
| Parser | Tokens | AST | Syntax validation, tree construction |
| Semantic Analyzer | AST | Typed AST | Type checking, scope analysis, inference |
| Codegen | Typed AST | LLVM IR | IR generation, optimization, linking |

## Detailed Pipeline

### 1. Lexer (Tokenization)

The lexer scans source code character by character and groups them into tokens. It handles:
- Keywords (`fn`, `if`, `struct`, etc.)
- Identifiers
- Operators (`+`, `-`, `*`, `/`, etc.)
- Literals (integers, floats, strings, booleans)
- Punctuation (`{}`, `()`, `;`, `,`, etc.)

The lexer tracks position information (filename, line, column) for error reporting.

### 2. Parser (AST Generation)

The parser uses recursive descent to build an Abstract Syntax Tree from tokens. It validates:
- Correct keyword ordering
- Matching parentheses/braces
- Proper expression syntax

The parser produces a tree of AST nodes representing program structure without type information.

### 3. Semantic Analyzer (Type Checking)

The semantic analyzer performs:
- **Scope Analysis**: Tracks variable declarations across scopes
- **Type Checking**: Validates type compatibility
- **Type Inference**: Determines types for untyped expressions
- **Operator Validation**: Ensures operators are used with valid operand types

This stage attaches type information to AST nodes, producing a typed AST.

### 4. Code Generation (LLVM IR)

The code generator:
- Maps Telia types to LLVM types
- Generates LLVM IR instructions
- Emits function bodies, control flow, expressions

After IR generation:
1. LLVM optimizer optimizes the IR
2. Clang links and produces final executable

## Data Flow

```
┌─────────────────────────────────────────────────────────────────────┐
│                         Source Code (.t)                           │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  Lexer                                                            │
│  - Scans characters                                               │
│  - Produces tokens with position info                             │
│  - Keywords, identifiers, literals, operators, punctuation         │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  Parser                                                           │
│  - Recursive descent parsing                                       │
│  - Builds AST without type info                                   │
│  - Validates syntax structure                                      │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  Semantic Analyzer                                                │
│  - Resolves types (type checking, inference)                      │
│  - Validates operators                                            │
│  - Manages scope (variable lookup)                                │
│  - Attaches type information to AST nodes                        │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  Code Generator                                                   │
│  - Maps types to LLVM types                                        │
│  - Emits LLVM IR                                                   │
│  - Creates function definitions, control flow, expressions        │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  LLVM Pipeline                                                    │
│  - opt -O3: Optimization passes                                   │
│  - clang: Linking with libc                                       │
│  - Final executable                                               │
└─────────────────────────────────────────────────────────────────────┘
```

## Error Handling

Errors are collected throughout the pipeline:
- Lexer: Invalid characters, malformed tokens
- Parser: Unexpected tokens, syntax errors
- Semantic Analyzer: Type mismatches, undefined variables, scope errors
- Codegen: IR generation failures

The CLI stops compilation if any stage reports errors.

## Extension Points

When adding new language features, each stage typically needs updates:

1. **Lexer**: Add new token kinds if new keywords/syntax needed
2. **Parser**: Add new AST node kinds and parsing logic
3. **Semantic Analyzer**: Add type checking and inference rules
4. **Codegen**: Add LLVM IR generation for new constructs
