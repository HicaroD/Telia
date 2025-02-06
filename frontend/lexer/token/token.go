package token

type Token struct {
	Lexeme []byte
	Kind   Kind
	Pos    Pos
}

func New(lexeme []byte, kind Kind, position Pos) *Token {
	return &Token{Lexeme: lexeme, Kind: kind, Pos: position}
}

func (token *Token) Name() string {
	if token.Kind == ID || token.Kind == STRING_LITERAL {
		return string(token.Lexeme)
	}
	return token.Kind.String()
}
