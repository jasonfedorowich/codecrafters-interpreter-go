package scanner

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type TokenizerSuite struct {
	suite.Suite
}

func (ts *TokenizerSuite) SetupTest() {

}

func (ts *TokenizerSuite) TearDownTest() {

}

func TestTokenizerSuite(t *testing.T) {
	suite.Run(t, new(TokenizerSuite))
}

func (ts *TokenizerSuite) TestTokenizer() {
	ts.Run("test with braces", func() {
		tokenizer := NewTokenizer("(()}}{")
		tokens, _ := tokenizer.Tokenize()
		ts.Assert().Equal(7, len(tokens))
		ts.Assert().Equal(EOF, tokens[6].Type)
		ts.Assert().Equal(LEFT_PAREN, tokens[0].Type)
		ts.Assert().Equal(LEFT_PAREN, tokens[1].Type)
		ts.Assert().Equal(RIGHT_PAREN, tokens[2].Type)
	})
	ts.Run("test with errors", func() {
		tokenizer := NewTokenizer("@")
		_, errCode := tokenizer.Tokenize()
		ts.Assert().Equal(65, errCode)
	})
	ts.Run("test with ==", func() {
		tokenizer := NewTokenizer("===")
		tokens, _ := tokenizer.Tokenize()
		ts.Assert().Equal(EQUAL_EQUAL, tokens[0].Type)
		ts.Assert().Equal(EQUAL, tokens[1].Type)
	})
	ts.Run("test with comment no line", func() {
		tokenizer := NewTokenizer("()// hello world")
		tokens, _ := tokenizer.Tokenize()
		ts.Assert().Len(tokens, 3)
	})
	ts.Run("test with slash", func() {
		tokenizer := NewTokenizer("/()")
		tokens, _ := tokenizer.Tokenize()
		ts.Assert().Len(tokens, 4)
	})
	ts.Run("test with new line comments", func() {
		tokenizer := NewTokenizer("/()///hello world\n(")
		tokens, _ := tokenizer.Tokenize()
		ts.Assert().Len(tokens, 5)
	})
	ts.Run("test with string", func() {
		tokenizer := NewTokenizer("\"hello world\"")
		tokens, _ := tokenizer.Tokenize()
		ts.Assert().Equal("hello world", tokens[0].Literal.s)
	})
	ts.Run("test with invalid string", func() {
		tokenizer := NewTokenizer("\"hello world")
		_, err := tokenizer.Tokenize()
		ts.Assert().Equal(65, err)
	})
	ts.Run("test with numbers", func() {
		tokenizer := NewTokenizer("12345")
		tokens, _ := tokenizer.Tokenize()
		ts.Assert().Equal(float64(12345), tokens[0].Literal.s)
	})
	ts.Run("test number with period", func() {
		tokenizer := NewTokenizer("12345.0901")
		tokens, _ := tokenizer.Tokenize()
		ts.Assert().Equal(12345.0901, tokens[0].Literal.s)
	})
}
