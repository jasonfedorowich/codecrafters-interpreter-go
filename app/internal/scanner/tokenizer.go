package scanner

type Tokenizer struct {
	p      int
	tokens []*Token
	source []rune
}

func NewTokenizer(source string) *Tokenizer {
	return &Tokenizer{
		source: []rune(source),
	}
}

func (t *Tokenizer) Tokenize() []*Token {
	for !t.isEnd() {
		t.scanToken()
	}
	t.addToken(EOF, "", nil)

	return t.tokens
}

func (t *Tokenizer) scanToken() {
	r := t.advance()
	switch r {
	case '(':
		t.addToken(LEFT_PAREN, "(", nil)
	case ')':
		t.addToken(RIGHT_PAREN, ")", nil)
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
