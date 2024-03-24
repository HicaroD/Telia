package kind

import "log"

// TODO: define all token kinds here

type TokenKind int

const (
	Id TokenKind = iota
)

func (kind TokenKind) String() string {
	switch kind {
	case Id:
		return "Id"
	default:
		log.Fatalf("Invalid token kind: %d", kind)
	}
	return ""
}
