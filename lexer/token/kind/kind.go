package kind

import (
	"log"
	"strconv"
)

type TokenKind int

const (
	// EOF
	EOF TokenKind = iota
	INVALID

	// Identifier
	ID

	// Literals
	INTEGER_LITERAL
	STRING_LITERAL
	TRUE_BOOL_LITERAL
	FALSE_BOOL_LITERAL

	// Keywords
	FN
	RETURN
	EXTERN
	IF
	ELIF
	ELSE
	NOT
	AND
	OR

	// Types
	BOOL_TYPE // bool

	INT_TYPE // int
	I8_TYPE  // i8
	I16_TYPE // i16
	I32_TYPE // i32
	I64_TYPE // i64

	UINT_TYPE // uint
	U8_TYPE   // u8
	U16_TYPE  // u16
	U32_TYPE  // u32
	U64_TYPE  // u64

	// This type is not explicit. We don't have a keyword for this, the absence
	// of an explicit type means a void type
	VOID_TYPE

	// (
	OPEN_PAREN
	// )
	CLOSE_PAREN

	// {
	OPEN_CURLY
	// }
	CLOSE_CURLY

	// ,
	COMMA

	// ;
	SEMICOLON

	// .
	DOT
	// ..
	DOT_DOT
	// ...
	DOT_DOT_DOT

	// =
	EQUAL
	// :=
	COLON_EQUAL
	// !=
	BANG_EQUAL
	// ==
	EQUAL_EQUAL

	// >
	GREATER
	// >=
	GREATER_EQ
	// <
	LESS
	// <=
	LESS_EQ

	// +
	PLUS
	// -
	MINUS
	// *
	STAR
	// /
	SLASH
)

var KEYWORDS map[string]TokenKind = map[string]TokenKind{
	"fn":     FN,
	"return": RETURN,
	"extern": EXTERN,
	"if":     IF,
	"elif":   ELIF,
	"else":   ELSE,
	"not":    NOT,
	"and":    AND,
	"or":     OR,

	"true":  TRUE_BOOL_LITERAL,
	"false": FALSE_BOOL_LITERAL,

	"bool": BOOL_TYPE,

	"int": INT_TYPE,
	"i8":  I8_TYPE,
	"i16": I16_TYPE,
	"i32": I32_TYPE,
	"i64": I64_TYPE,

	"uint": UINT_TYPE,
	"u8":   U8_TYPE,
	"u16":  U16_TYPE,
	"u32":  U32_TYPE,
	"u64":  U64_TYPE,
}

var BASIC_TYPES map[TokenKind]bool = map[TokenKind]bool{
	VOID_TYPE: true,
	BOOL_TYPE: true,
	INT_TYPE:  true,
	I8_TYPE:   true,
	I16_TYPE:  true,
	I32_TYPE:  true,
	I64_TYPE:  true,
	UINT_TYPE: true,
	U8_TYPE:   true,
	U16_TYPE:  true,
	U32_TYPE:  true,
	U64_TYPE:  true,
}

var LITERAL_KIND map[TokenKind]bool = map[TokenKind]bool{
	INTEGER_LITERAL:          true,
	STRING_LITERAL:           true,
	TRUE_BOOL_LITERAL:        true,
	FALSE_BOOL_LITERAL:       true,
}

var NUMERIC_TYPES map[TokenKind]bool = map[TokenKind]bool{
	INT_TYPE:  true,
	I8_TYPE:   true,
	I16_TYPE:  true,
	I32_TYPE:  true,
	I64_TYPE:  true,
	UINT_TYPE: true,
	U8_TYPE:   true,
	U16_TYPE:  true,
	U32_TYPE:  true,
	U64_TYPE:  true,
}

var LOGICAL_OP map[TokenKind]bool = map[TokenKind]bool{
	AND:         true,
	OR:          true,
	BANG_EQUAL:  true,
	EQUAL_EQUAL: true,
	GREATER:     true,
	GREATER_EQ:  true,
	LESS:        true,
	LESS_EQ:     true,
}

func (kind TokenKind) BitSize() int {
	switch kind {
	case INT_TYPE, UINT_TYPE:
		return strconv.IntSize
	case BOOL_TYPE:
		return 1
	case I8_TYPE, U8_TYPE:
		return 8
	case I16_TYPE, U16_TYPE:
		return 16
	case I32_TYPE, U32_TYPE:
		return 32
	case I64_TYPE, U64_TYPE:
		return 64
	default:
		return -1
	}
}

func (kind TokenKind) String() string {
	switch kind {
	case EOF:
		return "end of file"
	case INVALID:
		return "INVALID"
	case ID:
		return "identifier"
	case INTEGER_LITERAL:
		return "integer literal"
	case STRING_LITERAL:
		return "string literal"
	case TRUE_BOOL_LITERAL:
		return "true"
	case FALSE_BOOL_LITERAL:
		return "false"
	case FN:
		return "fn"
	case RETURN:
		return "return"
	case EXTERN:
		return "extern"
	case IF:
		return "if"
	case ELIF:
		return "elif"
	case ELSE:
		return "else"
	case NOT:
		return "not"
	case AND:
		return "and"
	case OR:
		return "or"
	case BOOL_TYPE:
		return "bool"
	case INT_TYPE:
		return "int"
	case I8_TYPE:
		return "i8"
	case I16_TYPE:
		return "i16"
	case I32_TYPE:
		return "i32"
	case I64_TYPE:
		return "i64"
	case UINT_TYPE:
		return "uint"
	case U8_TYPE:
		return "u8"
	case U16_TYPE:
		return "u16"
	case U32_TYPE:
		return "u32"
	case U64_TYPE:
		return "u64"
	case VOID_TYPE:
		return "void"
	case OPEN_PAREN:
		return "("
	case CLOSE_PAREN:
		return ")"
	case OPEN_CURLY:
		return "{"
	case CLOSE_CURLY:
		return "}"
	case COMMA:
		return ","
	case SEMICOLON:
		return ";"
	case DOT:
		return "."
	case DOT_DOT:
		return ".."
	case DOT_DOT_DOT:
		return "..."
	case EQUAL:
		return "="
	case COLON_EQUAL:
		return ":="
	case BANG_EQUAL:
		return "!="
	case EQUAL_EQUAL:
		return "=="
	case GREATER:
		return ">"
	case GREATER_EQ:
		return ">="
	case LESS:
		return "<"
	case LESS_EQ:
		return "<="
	case PLUS:
		return "+"
	case MINUS:
		return "-"
	case STAR:
		return "*"
	case SLASH:
		return "/"
	default:
		log.Fatalf("String() method not defined for the following token kind '%d'", kind)
	}
	return ""
}
