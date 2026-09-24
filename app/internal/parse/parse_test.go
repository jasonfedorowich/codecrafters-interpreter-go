package parse

import (
	"testing"

	"github.com/codecrafters-io/interpreter-starter-go/app/internal/scanner"
	"github.com/stretchr/testify/suite"
)

type ParseSuite struct {
	suite.Suite
}

func (ps *ParseSuite) SetupTest() {

}

func (ps *ParseSuite) TearDownTest() {

}

func TestParseSuite(t *testing.T) {
	suite.Run(t, new(ParseSuite))
}

func (ps *ParseSuite) TestParser() {
	ps.Run("test with boolean", func() {
		tokenizer := scanner.NewTokenizer("true")
		tokens, _ := tokenizer.Tokenize()
		parser := NewParser(tokens)
		expr, err := parser.Parse()
		ps.Assert().IsType(&LiteralExpr{}, expr)
		ps.Equal(0, err)
	})
	ps.Run("test with nil", func() {
		tokenizer := scanner.NewTokenizer("nil")
		tokens, _ := tokenizer.Tokenize()
		parser := NewParser(tokens)
		expr, err := parser.Parse()
		ps.Assert().IsType(&LiteralExpr{}, expr)
		ps.Equal(0, err)
	})
	ps.Run("test with string", func() {
		tokenizer := scanner.NewTokenizer("\"hello world\"")
		tokens, _ := tokenizer.Tokenize()
		parser := NewParser(tokens)
		expr, err := parser.Parse()
		ps.Assert().IsType(&LiteralExpr{}, expr)
		ps.Equal(0, err)
	})
	ps.Run("test with number", func() {
		tokenizer := scanner.NewTokenizer("42.37")
		tokens, _ := tokenizer.Tokenize()
		parser := NewParser(tokens)
		expr, err := parser.Parse()
		ps.Assert().IsType(&LiteralExpr{}, expr)
		ps.Equal(0, err)
	})
	ps.Run("test with parenthesized expression", func() {
		tokenizer := scanner.NewTokenizer("(true)")
		tokens, _ := tokenizer.Tokenize()
		parser := NewParser(tokens)
		expr, err := parser.Parse()
		ps.Assert().IsType(&GroupingExpr{}, expr)
		ps.Equal(0, err)
	})
	ps.Run("test with unary expression", func() {
		tokenizer := scanner.NewTokenizer("!true")
		tokens, _ := tokenizer.Tokenize()
		parser := NewParser(tokens)
		expr, err := parser.Parse()
		ps.Assert().IsType(&UnaryExpr{}, expr)
		ps.Equal(0, err)
	})
	ps.Run("test with binary expression", func() {
		tokenizer := scanner.NewTokenizer("1+2")
		tokens, _ := tokenizer.Tokenize()
		parser := NewParser(tokens)
		expr, err := parser.Parse()
		ps.Assert().IsType(&BinaryExpr{}, expr)
		ps.Equal(0, err)
	})
	ps.Run("test with comparison", func() {
		tokenizer := scanner.NewTokenizer("1<2")
		tokens, _ := tokenizer.Tokenize()
		parser := NewParser(tokens)
		expr, err := parser.Parse()
		ps.Assert().IsType(&BinaryExpr{}, expr)
		ps.Equal(0, err)
	})
	ps.Run("test with equality", func() {
		tokenizer := scanner.NewTokenizer("1==2")
		tokens, _ := tokenizer.Tokenize()
		parser := NewParser(tokens)
		expr, err := parser.Parse()
		ps.Assert().IsType(&BinaryExpr{}, expr)
		ps.Equal(0, err)
	})
}
