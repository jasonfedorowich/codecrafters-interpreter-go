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
	LEFT_PAREN    TokenType = "LEFT_PAREN"
	RIGHT_PAREN   TokenType = "RIGHT_PAREN"
	LEFT_BRACE    TokenType = "LEFT_BRACE"
	RIGHT_BRACE   TokenType = "RIGHT_BRACE"
	COMMA         TokenType = "COMMA"
	DOT           TokenType = "DOT"
	STAR          TokenType = "STAR"
	PLUS          TokenType = "PLUS"
	MINUS         TokenType = "MINUS"
	SEMICOLON     TokenType = "SEMICOLON"
	EQUAL         TokenType = "EQUAL"
	EQUAL_EQUAL   TokenType = "EQUAL_EQUAL"
	BANG          TokenType = "BANG"
	BANG_EQUAL    TokenType = "BANG_EQUAL"
	LESS          TokenType = "LESS"
	LESS_EQUAL    TokenType = "LESS_EQUAL"
	GREATER       TokenType = "GREATER"
	GREATER_EQUAL TokenType = "GREATER_EQUAL"
	SLASH         TokenType = "SLASH"
	EOF           TokenType = "EOF"
)

type Token struct {
	Type    TokenType
	Lexeme  nullString
	Literal nullString
	Line    int
}

func NewToken(t TokenType, lexeme, literal *string, line int) *Token {
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
