package scanner

import "fmt"

type Tokenizer struct {
	p      int
	tokens []*Token
	errs   []error
	source []rune
}

func NewTokenizer(source string) *Tokenizer {
	return &Tokenizer{
		source: []rune(source),
	}
}

func (t *Tokenizer) Tokenize() ([]*Token, []error) {
	for !t.isEnd() {
		t.scanToken()
	}
	t.addToken(EOF, "", nil)

	return t.tokens, t.errs
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
		err := fmt.Errorf("[line %d] Error: Unexpected character: %c", 1, r)
		t.errs = append(t.errs, err)
		fmt.Println(err)
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
