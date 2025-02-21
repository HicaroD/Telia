package token

import (
	"log"
	"strconv"
)

type Kind int

const (
	// EOF
	EOF Kind = iota
	INVALID

	// Identifier
	ID

	LITERAL_START // literal start delimiter

	INTEGER_LITERAL
	STRING_LITERAL
	TRUE_BOOL_LITERAL
	FALSE_BOOL_LITERAL

	LITERAL_END // literal end delimiter

	FN
	FOR
	WHILE
	RETURN
	EXTERN
	PACKAGE
	IF
	ELIF
	ELSE
	NOT
	TYPE
	USE

	// Types
	BASIC_TYPE_START // basic type start delimiter

	BOOL_TYPE // bool

	NUMERIC_TYPE_START // numeric type start delimiter

	INTEGER_TYPE_START // integer type start delimiter
	UNTYPED_INT        // int
	I8_TYPE            // i8
	I16_TYPE           // i16
	I32_TYPE           // i32
	I64_TYPE           // i64
	UNTYPED_UINT       // uint
	U8_TYPE            // u8
	U16_TYPE           // u16
	U32_TYPE           // u32
	U64_TYPE           // u64
	INTEGER_TYPE_END   // integer type end delimiter

	NUMERIC_TYPE_END // numeric type end delimiter

	STRING_TYPE  // string
	CSTRING_TYPE // cstring

	UNTYPED_STRING // string

	// This type is not explicit. We don't have a keyword for this, the absence
	// of an explicit type means a void type
	VOID_TYPE

	BASIC_TYPE_END // basic type end delimiter

	OPEN_PAREN
	CLOSE_PAREN

	OPEN_CURLY
	CLOSE_CURLY

	OPEN_BRACKET
	CLOSE_BRACKET

	COMMA
	SEMICOLON

	DOT
	DOT_DOT
	DOT_DOT_DOT
	EQUAL
	COLON_EQUAL

	LOGICAL_OP_START // logical op start delimiter

	AND
	OR
	BANG_EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQ
	LESS
	LESS_EQ

	LOGICAL_OP_END // logical op end delimiter

	PLUS
	MINUS
	STAR
	SLASH

	SHARP
)

var KEYWORDS map[string]Kind = map[string]Kind{
	"fn":      FN,
	"for":     FOR,
	"while":   WHILE,
	"return":  RETURN,
	"extern":  EXTERN,
	"package": PACKAGE,
	"if":      IF,
	"elif":    ELIF,
	"else":    ELSE,
	"not":     NOT,
	"and":     AND,
	"or":      OR,
	"type":    TYPE,
	"use":     USE,

	"true":  TRUE_BOOL_LITERAL,
	"false": FALSE_BOOL_LITERAL,

	"bool": BOOL_TYPE,

	"int": UNTYPED_INT,
	"i8":  I8_TYPE,
	"i16": I16_TYPE,
	"i32": I32_TYPE,
	"i64": I64_TYPE,

	"uint":    UNTYPED_UINT,
	"u8":      U8_TYPE,
	"u16":     U16_TYPE,
	"u32":     U32_TYPE,
	"u64":     U64_TYPE,
	"string":  STRING_TYPE,
	"cstring": CSTRING_TYPE,
}

func (k Kind) BitSize() int {
	switch k {
	case UNTYPED_INT, UNTYPED_UINT:
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

func (k Kind) IsBasicType() bool {
	return k > BASIC_TYPE_START && k < BASIC_TYPE_END
}

func (k Kind) IsLiteral() bool {
	return k > LITERAL_START && k < LITERAL_END
}

func (k Kind) IsLogicalOp() bool {
	return k > LOGICAL_OP_START && k < LOGICAL_OP_END
}

func (k Kind) String() string {
	switch k {
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
	case FOR:
		return "for"
	case WHILE:
		return "while"
	case RETURN:
		return "return"
	case EXTERN:
		return "extern"
	case PACKAGE:
		return "package"
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
	case TYPE:
		return "type"
	case USE:
		return "use"
	case BOOL_TYPE:
		return "bool"
	case UNTYPED_INT:
		return "int"
	case I8_TYPE:
		return "i8"
	case I16_TYPE:
		return "i16"
	case I32_TYPE:
		return "i32"
	case I64_TYPE:
		return "i64"
	case UNTYPED_UINT:
		return "uint"
	case U8_TYPE:
		return "u8"
	case U16_TYPE:
		return "u16"
	case U32_TYPE:
		return "u32"
	case U64_TYPE:
		return "u64"
	case UNTYPED_STRING:
		return "untyped string"
	case STRING_TYPE:
		return "string"
	case CSTRING_TYPE:
		return "cstring"
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
	case OPEN_BRACKET:
		return "["
	case CLOSE_BRACKET:
		return "]"
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
	case SHARP:
		return "#"
	default:
		log.Fatalf("String() method not defined for the following token kind '%d'", k)
	}
	return ""
}
