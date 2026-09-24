package scanner

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/codecrafters-io/interpreter-starter-go/app/internal/error_code"
)

type Tokenizer struct {
	p      int
	tokens []*Token
	source []rune
	line   int
	eh     *error_code.ErrorHandler
}

func NewTokenizer(source string) *Tokenizer {
	return &Tokenizer{
		source: []rune(source),
		line:   1,
		eh:     &error_code.ErrorHandler{},
	}
}

func (t *Tokenizer) Tokenize() ([]*Token, int) {
	for !t.isEnd() {
		t.scanToken()
	}
	t.addToken(EOF, "", nil)

	return t.tokens, t.eh.ErrorCode
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
		if unicode.IsDigit(r) {
			t.addNumber()
			return
		}
		if unicode.IsLetter(r) || r == '_' {
			t.addIdentifier()
			return
		}
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

func (t *Tokenizer) addToken(tokenType TokenType, lexeme string, literal any) {
	t.tokens = append(t.tokens, NewToken(tokenType, lexeme, literal, t.line))
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
	t.addToken(STRING, fmt.Sprintf("\"%s\"", str), str)
}

func (t *Tokenizer) addNumber() {
	start := t.p - 1
	for !t.isEnd() && t.isNumber() {
		t.advance()
	}
	if !t.isEnd() && t.match('.') {
		t.advance()
		for !t.isEnd() && t.isNumber() {
			t.advance()
		}
	}
	str := string(t.source[start:t.p])
	number, err := strconv.ParseFloat(str, 64)
	if err != nil {
		t.error("[line %d] Error: Invalid number.", t.line)
		return
	}
	t.addToken(NUMBER, str, number)
}

func (t *Tokenizer) addIdentifier() {
	start := t.p - 1
	for !t.isEnd() && t.isIdentifier() {
		t.advance()
	}
	str := string(t.source[start:t.p])
	token, ok := reserved[str]
	if !ok {
		t.addToken(IDENTIFIER, str, nil)
		return
	}
	t.addToken(token, str, nil)

}

func (t *Tokenizer) isNumber() bool {
	return unicode.IsDigit(t.source[t.p])
}

func (t *Tokenizer) error(format string, a ...any) {
	t.eh.Error(65, format, a...)
}

func (t *Tokenizer) isIdentifier() bool {
	return unicode.IsLetter(t.source[t.p]) || t.source[t.p] == '_' || unicode.IsNumber(t.source[t.p])
}
