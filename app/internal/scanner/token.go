package scanner

import (
	"fmt"
	"reflect"
)

type TokenType string

type null struct {
	s any
}

func (ns null) String() string {
	if ns.s == nil {
		return "null"
	}
	if v := reflect.ValueOf(ns.s); v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return "null"
		}
		return fmt.Sprint(v.Elem().Interface())
	}
	return fmt.Sprint(ns.s)
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
	STRING        TokenType = "STRING"
	NUMBER        TokenType = "NUMBER"
	EOF           TokenType = "EOF"
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal null
	Line    int
}

func NewToken(t TokenType, lexeme string, literal any, line int) *Token {
	token := &Token{
		Type:   t,
		Line:   line,
		Lexeme: lexeme,
	}
	token.Literal.s = literal
	return token
}

func (t *Token) String() string {
	return fmt.Sprintf("%s %s %s", t.Type, t.Lexeme, t.Literal)
}
