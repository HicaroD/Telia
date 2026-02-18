# Type System

Telia has a rich type system with basic types, pointers, structs, and tuples. The type system supports type checking, inference, and compatibility rules.

## Type Kinds

### Expression Type Kinds

| Kind | Description |
|------|-------------|
| `Basic` | Basic types (integers, floats, bool, string) |
| `Id` | Named type (type alias or unresolved type) |
| `Struct` | Struct type |
| `Pointer` | Pointer type |
| `Alias` | Type alias |
| `Tuple` | Tuple type |

### Basic Types

Telia provides these basic types:

| Type | Description | Size |
|------|-------------|------|
| `bool` | Boolean | 1 bit |
| `i8`, `i16`, `i32`, `i64`, `i128` | Signed integers | 8-128 bits |
| `u8`, `u16`, `u32`, `u64`, `u128` | Unsigned integers | 8-128 bits |
| `int`, `uint` | Platform-sized integers | 32 or 64 bits |
| `f32`, `f64` | Floating point | 32 or 64 bits |
| `float` | Platform-sized float | 32 or 64 bits |
| `string` | String type | Pointer + length |
| `cstring` | C-compatible string | Pointer |
| `rawptr` | Raw pointer | Pointer |
| `void` | No type | N/A |

### Compound Types

#### Pointer Types

Pointers are created with `*Type` syntax:
- `*i32` - pointer to 32-bit integer
- `*User` - pointer to struct

Each pointer type tracks:
- The type it points to (pointee type)
- Whether the pointer was explicitly declared

#### Struct Types

Structs are defined with the `struct` keyword and contain named fields. Each struct type:
- References its declaration
- Has a set of fields with names and types

#### Tuple Types

Tuples are heterogeneous fixed-size collections:
- `(i32, string)` - tuple of integer and string
- `(i32, f64, bool)` - triple

Each tuple type tracks:
- Ordered list of element types

## Type Representation

### ExprType Structure

The expression type contains:
- **Kind**: The kind of type (basic, pointer, struct, etc.)
- **Underlying type**: The actual type struct (varies by kind)
- **Explicit flag**: Whether type was explicitly specified by user

The `Explicit` flag indicates whether the type was:
- **Explicit**: Type was written by user (e.g., `x: i32`)
- **Inferred**: Type was inferred from context (e.g., `x := 5`)

This distinction matters for type compatibility checking.

## Type Operations

### Type Equality

Two types are equal if they have the same kind and underlying representation:

- Basic types: Token kind must match
- Pointers: Pointee types must be equal
- Tuples: All element types must be equal
- Structs: Must reference same struct declaration

### Type Compatibility

Two types are compatible if they can be used in the same context:

- Integers: All integer types compatible with each other
- Floats: All float types compatible with each other
- Bool: Only compatible with bool
- String: Only compatible with string

### Type Promotion

Promotes untyped literals to explicit types:

- Untyped integers → `int`
- Untyped floats → `float`

## Operator Type Tables

### Unary Operators

Each unary operator defines:
- Valid operand types
- Result type (if applicable)
- Optional custom type resolution logic

Common unary operators:
- Numeric negation (`-`)
- Boolean negation (`not`)
- Address-of (`&`)
- Dereference (`*`)

### Binary Operators

Each binary operator defines:
- Valid operand types
- Result type (if applicable)
- Optional custom type resolution logic

Common binary operators:
- Arithmetic: `+`, `-`, `*`, `/`
- Comparison: `<`, `<=`, `>`, `>=`, `==`, `!=`
- Logical: `and`, `or`

## Type Checking in Semantic Analyzer

The semantic analyzer uses these type operations:

1. **Resolution**: Convert named types to actual type
2. **Matching**: Check if expression type matches expected type
3. **Inference**: Determine type from context when not explicit
4. **Validation**: Ensure operator operands are valid

## Adding New Types

When adding a new type:

1. **ExprTypeKind**:
   - Add new constant to the type kind enum

2. **Type Struct**:
   - Define struct for the new type

3. **ExprType Methods**:
   - Implement equality checking, numeric checks, pointer checks, etc.

4. **Semantic Analyzer**:
   - Add type inference and checking logic

5. **Code Generator**:
   - Add LLVM type mapping
