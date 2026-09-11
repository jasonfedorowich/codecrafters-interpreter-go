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
		tokenizer := NewTokenizer("<")
		_, errCode := tokenizer.Tokenize()
		ts.Assert().Equal(65, errCode)
	})
}
