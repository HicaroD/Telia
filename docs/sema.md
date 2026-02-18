# Semantic Analyzer

The semantic analyzer (sema) performs type checking, type inference, scope analysis, and validates that the program follows semantic rules beyond what the parser can enforce. It transforms the untyped AST into a typed AST with resolved symbols.

## Responsibilities

### 1. Scope Management

Tracks variable declarations across nested scopes and resolves identifier references:

- **Insert**: Add new symbol to current scope
- **LookupCurrentScope**: Find symbol in current scope only
- **LookupAcrossScopes**: Find symbol, traversing parent scopes

Each scope contains:
- Reference to parent scope (null for global)
- Map of symbol names to their definitions

### 2. Type Checking

Validates type compatibility in various contexts:

- **Variable declarations**: Expression type must match variable type
- **Assignments**: Right-hand side must be assignable to left-hand side
- **Function returns**: Return value must match function's return type
- **Function arguments**: Argument types must match parameter types

### 3. Type Inference

Determines types for expressions without explicit type annotations:

- **Variable initialization**: `x := 5` infers `x` as integer
- **Literal types**: Untyped literals get explicit types from context
- **Binary operations**: Result type derived from operand types

### 4. Operator Validation

Ensures operators are used with valid operand types via operator tables. Each operator defines:
- Which types it works with
- What result type it produces
- Optional custom type resolution logic

## Key Analysis Functions

### Package Checking

- Validates package name is not "main" (reserved)
- Checks all files in package
- Ensures main package has main function

### Function Checking

- Validates function attributes
- Resolves parameter and return types
- Checks function body
- Validates all paths return a value (for non-void functions)

### Block Checking

- Processes each statement in block
- Maintains scope chain
- Validates return type compatibility

### Expression Type Inference

The semantic analyzer has two modes:

1. **With Context**: Type is known (e.g., parameter matching)
   - Check if expression can conform to expected type
   - Apply type conversions if needed
   - Return the expected type (possibly modified)

2. **Without Context**: Type must be inferred (e.g., variable initialization)
   - Analyze expression based on its kind
   - Return inferred type and whether it's explicit

## Type Inference Process

### With Context

When the expected type is known:

1. Check if expression can conform to expected type
2. Apply type conversions if needed
3. Return the expected type (possibly modified)

### Without Context

When type must be inferred:

1. Analyze expression based on its kind
2. Return inferred type and whether it's explicit

## Checking Patterns

### Variable Declaration

1. Check if variable already exists in scope (for declarations)
2. Check expression type
3. Match expression type with variable type (or infer)
4. Insert into scope

### Function Call

1. Look up function in scope
2. Validate it's callable (function or prototype)
3. Check argument count matches parameter count
4. Check each argument type matches parameter type

### Control Flow

- Condition expressions must evaluate to boolean
- Each branch is type-checked with appropriate return type context
- Validates all paths return for non-void functions

## Error Handling

Errors are collected via diagnostics:

- **Undefined variable**: Variable not found in any scope
- **Redeclaration**: Variable already exists in current scope
- **Type mismatch**: Expression type doesn't match expected type
- **Invalid operator**: Operator used with invalid operand types
- **Not callable**: Attempting to call non-function value

## Adding New Constructs

When adding new language features:

1. **Type Inference**:
   - Add case for inferring type with context
   - Add case for inferring type without context

2. **Statement Checking**:
   - Add case to statement checking logic
   - Implement type validation

3. **Operator Support**:
   - Add entry to unary or binary operator tables
   - Define valid types and result type
