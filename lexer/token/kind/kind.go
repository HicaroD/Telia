package kind

import "log"

// TODO: define all token kinds here

type TokenKind int

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

	// (
	OPEN_PAREN
	// )
	CLOSE_PAREN

	// {
	OPEN_CURLY
	// }
	CLOSE_CURLY

	// ;
	SEMICOLON
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
		return "INTEGER_LITERAL"
	case FN:
		return "FN"
	case RETURN:
		return "RETURN"
	case OPEN_PAREN:
		return "("
	case CLOSE_PAREN:
		return ")"
	case OPEN_CURLY:
		return "{"
	case CLOSE_CURLY:
		return "}"
	case SEMICOLON:
		return ";"
	default:
		log.Fatalf("String() method not defined for the following token kind '%d'", kind)
	}
	return ""
}
