package token

type Token struct {
	Lexeme   string
	Kind     Kind
	Position Position
}

func New(lexeme string, kind Kind, position Position) *Token {
	return &Token{Lexeme: lexeme, Kind: kind, Position: position}
}

func (token *Token) Name() string {
	if token.Kind == ID {
		return token.Lexeme
	}
	return token.Kind.String()
}
