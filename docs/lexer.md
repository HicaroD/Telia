# Lexer (Tokenizer)

The lexer transforms raw source code text into a stream of tokens. It is the first stage of compilation and operates at the character level.

## Token Structure

Each token contains:
- **Kind**: The type of token (keyword, operator, identifier, etc.)
- **Lexeme**: The actual text that was scanned
- **Position**: Source location (filename, line, column)

## Token Categories

### Keywords

Reserved words with special meaning:

| Token | Description |
|-------|-------------|
| `fn` | Function declaration |
| `struct` | Struct declaration |
| `if`, `elif`, `else` | Conditional branching |
| `for`, `while` | Loop constructs |
| `return` | Return statement |
| `extern` | External function declarations |
| `use` | Import statement |
| `defer` | Deferred execution |
| `type` | Type alias |
| `package` | Package declaration |

### Types

Built-in type keywords:

| Token | Description |
|-------|-------------|
| `bool` | Boolean type |
| `i8`, `i16`, `i32`, `i64`, `i128` | Signed integers |
| `u8`, `u16`, `u32`, `u64`, `u128` | Unsigned integers |
| `f32`, `f64` | Floating point |
| `int`, `uint` | Platform-sized integers |
| `float` | Platform-sized float |
| `string`, `cstring` | String types |
| `rawptr` | Raw pointer |
| `void` | No return type |

### Operators

Single and multi-character operators:

| Token | Description |
|-------|-------------|
| `+`, `-`, `*`, `/` | Arithmetic |
| `==`, `!=`, `<`, `<=`, `>`, `>=` | Comparison |
| `and`, `or`, `not` | Logical |
| `&` | Address-of |
| `=` | Assignment |
| `:=` | Declaration with inference |

### Punctuation

| Token | Description |
|-------|-------------|
| `(`, `)` | Parentheses |
| `{`, `}` | Braces |
| `[`, `]` | Brackets |
| `,`, `;` | Separators |
| `.`, `::` | Access operators |
| `:` | Type annotation |

## Lexer Implementation Patterns

### Character Scanning

The lexer scans character by character using a peek/advance pattern. It provides methods to:
- Peek at the next character without consuming it
- Skip/consume characters
- Check if the next character matches a specific value

### Token Classification

1. **Alphabetic start**: Could be keyword or identifier
   - Scan identifier characters
   - Check against keyword map
   - Return identifier or keyword token

2. **Digit start**: Numeric literal
   - Scan digits (base 10)
   - Check for decimal point for floats

3. **String delimiter**: String literal
   - Scan until closing quote
   - Handle escape sequences

4. **Operator characters**: Operator or punctuation
   - Check for multi-character operators (e.g., `==`, `!=`)
   - Otherwise return single character

### Position Tracking

The lexer maintains position information for error reporting:
- Filename being processed
- Current line number
- Current column number
- Directory containing source file

Line and column are incremented as newline and tab characters are encountered.

### Strict Newline Mode

The lexer supports a strict newline mode that treats unexpected newlines as errors. This is used in specific parsing contexts like struct literals where newlines would break the syntax.

## Adding New Tokens

When adding new syntax to the language:

1. **Token Kind**:
   - Add new constant to the token kind enum
   - Add keyword mapping if it's a reserved word

2. **Lexer**:
   - Add scanning logic for new token type
   - Handle in the main scan loop

3. **Parser** (if new syntax):
   - Add parsing logic for new token kind
   - Add to expression precedence if operator

## Error Handling

Lexer errors are collected via diagnostics:
- Invalid characters
- Unterminated string literals
- Malformed numeric literals

Errors include source position for accurate reporting.
