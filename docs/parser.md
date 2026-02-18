# Parser

The parser transforms a stream of tokens into an Abstract Syntax Tree (AST). It validates the syntactic structure of the program and produces a tree representation that captures the program's hierarchy.

## Parser Approach

The parser uses **recursive descent** with explicit precedence climbing for expressions. This means:
- Each syntactic construct has its own parsing function
- Expression precedence is handled through a chain of parsing functions
- The parser consumes tokens and builds node structures

## AST Node Kinds

### Declaration Nodes

| Kind | Description |
|------|-------------|
| `FnDecl` | Function declaration |
| `ExternDecl` | External function declarations |
| `StructDecl` | Struct definition |
| `TypeAliasDecl` | Type alias |
| `PkgDecl` | Package declaration |
| `UseDecl` | Import statement |

### Statement Nodes

| Kind | Description |
|------|-------------|
| `BlockStmt` | Block of statements |
| `VarStmt` | Variable declaration/assignment |
| `AssignmentStmt` | Assignment to existing variable |
| `ReturnStmt` | Return statement |
| `CondStmt` | If/elif/else statement |
| `ForLoopStmt` | For loop |
| `WhileLoopStmt` | While loop |
| `DeferStmt` | Deferred statement |

### Expression Nodes

| Kind | Description |
|------|-------------|
| `LiteralExpr` | Literal value (number, string, bool) |
| `IdExpr` | Identifier reference |
| `FnCall` | Function call |
| `BinaryExpr` | Binary operation |
| `UnaryExpr` | Unary operation |
| `TupleLiteralExpr` | Tuple literal |
| `StructExpr` | Struct literal |
| `FieldAccess` | Struct field access |
| `DerefPointerExpr` | Pointer dereference |
| `AddressOfExpr` | Address-of operation |
| `NamespaceAccess` | Namespace/module access |
| `NullPtrExpr` | Null pointer literal |

## Expression Precedence

The parser uses a precedence-climbing approach for expressions, from lowest to highest:

```
LogicalExpr      → ComparisonExpr (("and" | "or") ComparisonExpr)*
ComparisonExpr   → Term (("==" | "!=" | ">" | ">=" | "<" | "<=") Term)*
Term            → Factor (("+" | "-") Factor)*
Factor          → UnaryExpr (("*" | "/") UnaryExpr)*
UnaryExpr       → [("!" | "-" | "&" | "*")] PrimaryExpr
PrimaryExpr     → Literal | ID | "(" Expr ")" | FnCall | ...
```

Each level is parsed by a dedicated function that handles operators at that precedence level.

## Key Parsing Patterns

### Token Peeking

The parser peeks at upcoming tokens to decide which production to use. Common patterns include:
- Looking ahead to determine if we have a declaration vs statement
- Checking if the next token matches an expected kind
- Peering past parentheses/braces to understand expression structure

### Expect Pattern

For required tokens, the parser uses an expect pattern that:
- Checks if the next token matches the expected kind
- Returns an error with source position if not matched
- Advances past the token if matched

### Recursive Descent Structure

The parser organizes its functions hierarchically:

```
Program
  └─ Package
     └─ Declarations
        ├─ Function Declaration
        │   ├─ Parameters
        │   └─ Block
        │       └─ Statements
        │           ├─ If Statement
        │           ├─ For Loop
        │           ├─ While Loop
        │           ├─ Variable
        │           ├─ Return
        │           └─ Expression
        │               └─ (expression hierarchy by precedence)
        ├─ Struct Declaration
        ├─ Type Alias
        └─ Import
```

## Scope Management

The parser creates scopes for blocks and attaches them to AST nodes. Each scope:
- Has a parent scope reference for chained lookup
- Contains a map of declared symbols
- Is passed down during parsing and used by the semantic analyzer later

## Error Handling

Parser errors are reported via diagnostics:
- Unexpected token errors
- Missing expected tokens (parentheses, braces)
- Invalid expression syntax

Errors include source position for accurate reporting.

## Adding New Syntax

When adding new language features:

1. **AST Node**:
   - Add new node kind constant
   - Add corresponding struct type

2. **Parser**:
   - Add parsing function for new construct
   - Integrate into appropriate precedence level
   - Handle in statement or expression parsing as needed

3. **Semantic Analyzer**:
   - Add type checking for new node kind
   - Add scope handling if new scoping rules apply

4. **Code Generator**:
   - Add IR generation for new node kind
