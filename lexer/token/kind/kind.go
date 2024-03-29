package kind

import (
	"log"
)

type TokenKind int

var KEYWORDS map[string]TokenKind = map[string]TokenKind{
	"fn":     FN,
	"return": RETURN,
	"extern": EXTERN,

	"bool": BOOL_TYPE,
	"i8":   I8_TYPE,
	"i16":  I16_TYPE,
	"i32":  I32_TYPE,
	"i64":  I64_TYPE,
	"i128": I128_TYPE,
}

var BASIC_TYPES map[TokenKind]bool = map[TokenKind]bool{
	BOOL_TYPE: true,
	I8_TYPE:   true,
	I16_TYPE:  true,
	I32_TYPE:  true,
	I64_TYPE:  true,
	I128_TYPE: true,
}

const (
	// EOF
	EOF TokenKind = iota
	INVALID

	// Identifier, literals and keywords
	ID
	INTEGER_LITERAL
	STRING_LITERAL
	FN
	RETURN
	EXTERN

	// Types
	BOOL_TYPE
	I8_TYPE
	I16_TYPE
	I32_TYPE
	I64_TYPE
	I128_TYPE

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

	// ..
	DOT_DOT
	// ...
	DOT_DOT_DOT

	// *
	STAR
)

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
	case FN:
		return "fn"
	case RETURN:
		return "return"
	case EXTERN:
		return "extern"
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
	case I128_TYPE:
		return "i128"
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
	case DOT_DOT:
		return ".."
	case DOT_DOT_DOT:
		return "..."
	case STAR:
		return "*"
	default:
		log.Fatalf("String() method not defined for the following token kind '%d'", kind)
	}
	return ""
}
