package scanner

import "fmt"

type TokenType string

type nullString struct {
	s *string
}

func (ns nullString) String() string {
	if ns.s == nil {
		return "null"
	}
	return *ns.s
}

const (
	LEFT_PAREN  TokenType = "LEFT_PAREN"
	RIGHT_PAREN TokenType = "RIGHT_PAREN"
	LEFT_BRACE  TokenType = "LEFT_BRACE"
	RIGHT_BRACE TokenType = "RIGHT_BRACE"
	COMMA       TokenType = ","
	DOT         TokenType = "."
	STAR        TokenType = "*"
	PLUS        TokenType = "+"
	MINUS       TokenType = "-"
	SEMI_COLON  TokenType = ";"
	EOF         TokenType = "EOF"
)

type Token struct {
	Type    TokenType
	Lexeme  nullString
	Literal nullString
}

func NewToken(t TokenType, lexeme, literal *string) *Token {
	token := &Token{
		Type: t,
	}
	token.Lexeme.s = lexeme
	token.Literal.s = literal
	return token
}

func (t *Token) String() string {
	return fmt.Sprintf("%s %s %s", t.Type, t.Lexeme, t.Literal)
}
