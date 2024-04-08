package kind

import (
	"log"
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

	// Types
	BOOL_TYPE
	I8_TYPE
	I16_TYPE
	I32_TYPE
	I64_TYPE
	VOID_TYPE
	// I128_TYPE

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

	// *
	STAR

	// =
	EQUAL
	// ==
	EQUAL_EQUAL

	// :=
	COLON_EQUAL

	// -
	MINUS
)

var KEYWORDS map[string]TokenKind = map[string]TokenKind{
	"fn":     FN,
	"return": RETURN,
	"extern": EXTERN,
	"if":     IF,
	"elif":   ELIF,
	"else":   ELSE,

	"true":  TRUE_BOOL_LITERAL,
	"false": FALSE_BOOL_LITERAL,

	"bool": BOOL_TYPE,
	"i8":   I8_TYPE,
	"i16":  I16_TYPE,
	"i32":  I32_TYPE,
	"i64":  I64_TYPE,
	// "i128": I128_TYPE,
}

var BASIC_TYPES map[TokenKind]bool = map[TokenKind]bool{
	BOOL_TYPE: true,
	I8_TYPE:   true,
	I16_TYPE:  true,
	I32_TYPE:  true,
	I64_TYPE:  true,
	VOID_TYPE: true,
	// I128_TYPE: true,
}

func (kind TokenKind) String() string {
	switch kind {
	case EOF:
		return "EOF"
	case INVALID:
		return "INVALID"
	case ID:
		return "ID"
	case INTEGER_LITERAL:
		return "INTEGER_LITERAL"
	case STRING_LITERAL:
		return "STRING_LITERAL"
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
	case BOOL_TYPE:
		return "bool"
	case I8_TYPE:
		return "i8"
	case I16_TYPE:
		return "i16"
	case I32_TYPE:
		return "i32"
	case I64_TYPE:
		return "i64"
	case VOID_TYPE:
		return "VOID"
	// case I128_TYPE:
	// 	return "i128"
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
	case STAR:
		return "*"
	case EQUAL:
		return "="
	case EQUAL_EQUAL:
		return "=="
	case MINUS:
		return "-"
	case COLON_EQUAL:
		return ":="
	default:
		log.Fatalf("String() method not defined for the following token kind '%d'", kind)
	}
	return ""
}
