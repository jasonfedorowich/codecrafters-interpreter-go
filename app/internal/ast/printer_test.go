package ast

import (
	"testing"

	"github.com/codecrafters-io/interpreter-starter-go/app/internal/parse"
	"github.com/codecrafters-io/interpreter-starter-go/app/internal/scanner"
	"github.com/stretchr/testify/suite"
)

type PrintTestSuite struct {
	suite.Suite
	astPrinter *AstPrinter
}

func (ps *PrintTestSuite) SetupTest() {
	ps.astPrinter = &AstPrinter{}
}

func (ps *PrintTestSuite) TearDownTest() {

}

func TestPrintSuite(t *testing.T) {
	suite.Run(t, new(PrintTestSuite))
}

func (ps *PrintTestSuite) TestPrint() {
	ps.Run("test with boolean", func() {
		tokenizer := scanner.NewTokenizer("true")
		tokens, _ := tokenizer.Tokenize()
		parser := parse.NewParser(tokens)
		expr, err := parser.Parse()
		s := expr.Accept(ps.astPrinter)
		ps.Equal("true", s)
		ps.NoError(err)
	})
	ps.Run("test with nil", func() {
		tokenizer := scanner.NewTokenizer("nil")
		tokens, _ := tokenizer.Tokenize()
		parser := parse.NewParser(tokens)
		expr, err := parser.Parse()
		s := expr.Accept(ps.astPrinter)
		ps.Equal("nil", s)
		ps.NoError(err)
	})
	ps.Run("test with parenthesized expression", func() {
		tokenizer := scanner.NewTokenizer("(true)")
		tokens, _ := tokenizer.Tokenize()
		parser := parse.NewParser(tokens)
		expr, err := parser.Parse()
		s := expr.Accept(ps.astPrinter)
		ps.Equal("(group true)", s)
		ps.NoError(err)
	})
}
