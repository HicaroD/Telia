package token

import (
	"log"
)

type Kind int

const (
	// EOF
	EOF Kind = iota
	INVALID

	NEWLINE
	ID

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
	DEFER
	STRUCT

	// Types
	BASIC_TYPE_START // basic type start delimiter

	RAWPTR_TYPE // rawptr
	BOOL_TYPE   // bool

	NUMERIC_TYPE_START // numeric type start delimiter

	INTEGER_TYPE_START // integer type start delimiter
	INT_TYPE           // int
	I8_TYPE            // i8
	I16_TYPE           // i16
	I32_TYPE           // i32
	I64_TYPE           // i64
	I128_TYPE          // i128
	UINT_TYPE          // int
	U8_TYPE            // u8
	U16_TYPE           // u16
	U32_TYPE           // u32
	U64_TYPE           // u64
	U128_TYPE          // u128
	INTEGER_TYPE_END   // integer type end delimiter

	FLOAT_TYPE_START // float type start delimiter
	FLOAT_TYPE       // float
	F32_TYPE         // f32
	F64_TYPE         // f64
	FLOAT_TYPE_END   // float type end delimiter

	NUMERIC_TYPE_END // numeric type end delimiter

	STRING_TYPE  // string
	CSTRING_TYPE // cstring

	UNTYPED_START // literal start delimiter

	UNTYPED_STRING
	UNTYPED_FLOAT
	UNTYPED_INT
	UNTYPED_BOOL

	UNTYPED_END // literal end delimiter

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
	COLON
	COLON_EQUAL
	COLON_COLON

	LOGICAL_OP_START // logical op start delimiter

	AND
	OR
	CMP_OP_START // comparison op start delimiter
	BANG_EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQ
	LESS
	LESS_EQ
	CMP_OP_END // comparison op end delimiter

	LOGICAL_OP_END // logical op end delimiter

	PLUS
	MINUS
	STAR
	SLASH

	SHARP
	AT
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
	"defer":   DEFER,
	"struct":  STRUCT,

	"true":  UNTYPED_BOOL,
	"false": UNTYPED_BOOL,

	"rawptr": RAWPTR_TYPE,
	"bool":   BOOL_TYPE,

	"int":  INT_TYPE,
	"i8":   I8_TYPE,
	"i16":  I16_TYPE,
	"i32":  I32_TYPE,
	"i64":  I64_TYPE,
	"i128": I128_TYPE,

	"uint": UINT_TYPE,
	"u8":   U8_TYPE,
	"u16":  U16_TYPE,
	"u32":  U32_TYPE,
	"u64":  U64_TYPE,
	"u128": U128_TYPE,

	"float": FLOAT_TYPE,
	"f32":   F32_TYPE,
	"f64":   F64_TYPE,

	"string":  STRING_TYPE,
	"cstring": CSTRING_TYPE,
}

func (k Kind) BitSize() int {
	switch k {
	case BOOL_TYPE, UNTYPED_BOOL:
		return 1
	case I8_TYPE, U8_TYPE:
		return 8
	case I16_TYPE, U16_TYPE:
		return 16
	case I32_TYPE, U32_TYPE, INT_TYPE, UINT_TYPE, F32_TYPE, FLOAT_TYPE, UNTYPED_FLOAT, UNTYPED_INT:
		return 32
	case I64_TYPE, U64_TYPE, F64_TYPE:
		return 64
	case I128_TYPE, U128_TYPE:
		return 128
	default:
		return -1
	}
}

func (k Kind) IsBasicType() bool {
	return k > BASIC_TYPE_START && k < BASIC_TYPE_END
}

func (k Kind) IsUntyped() bool {
	return k > UNTYPED_START && k < UNTYPED_END
}

func (k Kind) IsNumeric() bool {
	return k > NUMERIC_TYPE_START && k < NUMERIC_TYPE_END || k == UNTYPED_INT || k == UNTYPED_FLOAT
}

func (k Kind) IsInteger() bool {
	return k > INTEGER_TYPE_START && k < INTEGER_TYPE_END || k == UNTYPED_INT
}

func (k Kind) IsFloat() bool {
	return k > FLOAT_TYPE_START && k < FLOAT_TYPE_END || k == UNTYPED_FLOAT
}

func (k Kind) IsStringLiteral() bool { return k == UNTYPED_STRING }
func (k Kind) IsStringType() bool {
	return k == STRING_TYPE || k == CSTRING_TYPE
}

func (k Kind) IsLogicalOp() bool {
	return k > LOGICAL_OP_START && k < LOGICAL_OP_END
}

func (k Kind) IsCmpOp() bool {
	return k > CMP_OP_START && k < CMP_OP_END
}

func (k Kind) String() string {
	switch k {
	case EOF:
		return "end of file"
	case INVALID:
		return "INVALID"
	case NEWLINE:
		return "new line"
	case ID:
		return "identifier"
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
	case DEFER:
		return "defer"
	case STRUCT:
		return "struct"
	case RAWPTR_TYPE:
		return "rawptr"
	case BOOL_TYPE:
		return "bool"
	case UNTYPED_FLOAT:
		return "untyped float"
	case UNTYPED_INT:
		return "untyped int"
	case UNTYPED_BOOL:
		return "untyped bool"
	case FLOAT_TYPE:
		return "f32"
	case F32_TYPE:
		return "f32"
	case F64_TYPE:
		return "f64"
	case INT_TYPE:
		return "int"
	case UINT_TYPE:
		return "uint"
	case I8_TYPE:
		return "i8"
	case I16_TYPE:
		return "i16"
	case I32_TYPE:
		return "i32"
	case I64_TYPE:
		return "i64"
	case I128_TYPE:
		return "i128"
	case U8_TYPE:
		return "u8"
	case U16_TYPE:
		return "u16"
	case U32_TYPE:
		return "u32"
	case U64_TYPE:
		return "u64"
	case U128_TYPE:
		return "u128"
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
	case COLON:
		return ":"
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
	case COLON_COLON:
		return "::"
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
	case AT:
		return "@"
	default:
		log.Fatalf("String() method not defined for the following token kind '%d'", k)
	}
	return ""
}
