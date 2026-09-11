package scanner

import (
	"fmt"
	"os"
)

type Tokenizer struct {
	p        int
	tokens   []*Token
	exitCode int
	source   []rune
}

func NewTokenizer(source string) *Tokenizer {
	return &Tokenizer{
		source: []rune(source),
	}
}

func (t *Tokenizer) Tokenize() ([]*Token, int) {
	for !t.isEnd() {
		t.scanToken()
	}
	t.addToken(EOF, "", nil)

	return t.tokens, t.exitCode
}

func (t *Tokenizer) scanToken() {
	r := t.advance()
	switch r {
	case '(':
		t.addToken(LEFT_PAREN, "(", nil)
	case ')':
		t.addToken(RIGHT_PAREN, ")", nil)
	case '{':
		t.addToken(LEFT_BRACE, "{", nil)
	case '}':
		t.addToken(RIGHT_BRACE, "}", nil)
	case '*':
		t.addToken(STAR, "*", nil)
	case '+':
		t.addToken(PLUS, "+", nil)
	case ',':
		t.addToken(COMMA, ",", nil)
	case '.':
		t.addToken(DOT, ".", nil)
	case '-':
		t.addToken(MINUS, "-", nil)
	case ';':
		t.addToken(SEMICOLON, ";", nil)
	default:
		t.exitCode = 65
		fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %c\n", 1, r)
	}
}

func (t *Tokenizer) isEnd() bool {
	return t.p >= len(t.source)
}

func (t *Tokenizer) advance() rune {
	r := t.source[t.p]
	t.p++
	return r
}

func (t *Tokenizer) addToken(tokenType TokenType, lexeme string, literal *string) {
	t.tokens = append(t.tokens, NewToken(tokenType, &lexeme, literal))
}
