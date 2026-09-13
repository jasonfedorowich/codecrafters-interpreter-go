package scanner

import (
	"fmt"
	"os"
	"strconv"
)

type Tokenizer struct {
	p        int
	tokens   []*Token
	exitCode int
	source   []rune
	line     int
}

func NewTokenizer(source string) *Tokenizer {
	return &Tokenizer{
		source: []rune(source),
		line:   1,
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
	case '"':
		t.addString()
	case '/':
		t.handleComment()
	case ' ', '\t', '\r':
		return
	case '\n':
		t.line++
	case '!':
		t.addTokenMatch('=', BANG_EQUAL, BANG, "!=", "!")
	case '=':
		t.addTokenMatch('=', EQUAL_EQUAL, EQUAL, "==", "=")
	case '<':
		t.addTokenMatch('=', LESS_EQUAL, LESS, "<=", "<")
	case '>':
		t.addTokenMatch('=', GREATER_EQUAL, GREATER, ">=", ">")
	default:
		t.error("[line %d] Error: Unexpected character: %c\n", t.line, r)
	}
}

func (t *Tokenizer) addTokenMatch(match rune, altToken, token TokenType, altLexeme, lexeme string) {
	if t.match(match) {
		t.advance()
		t.addToken(altToken, altLexeme, nil)
	} else {
		t.addToken(token, lexeme, nil)
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
	t.tokens = append(t.tokens, NewToken(tokenType, &lexeme, literal, t.line))
}

func (t *Tokenizer) match(r rune) bool {
	if t.isEnd() {
		return false
	}
	return t.source[t.p] == r
}

func (t *Tokenizer) handleComment() {
	if t.match('/') {
		for !t.isEnd() && !t.match('\n') {
			t.advance()
		}
	} else {
		t.addToken(SLASH, "/", nil)
	}
}

func (t *Tokenizer) addString() {
	start := t.p
	for !t.isEnd() && !t.match('"') {
		t.advance()
	}
	if t.isEnd() || !t.match('"') {
		t.error("[line %d] Error: Unterminated string.", t.line)
		return
	}
	str := string(t.source[start:t.p])
	t.advance()
	t.addToken(STRING, strconv.Quote(str), &str)
}

func (t *Tokenizer) error(format string, a ...any) {
	t.exitCode = 65
	fmt.Fprintf(os.Stderr, format, a...)
}
