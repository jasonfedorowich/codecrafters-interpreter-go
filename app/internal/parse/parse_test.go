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
		tokenizer := scanner.NewTokenizer("true;")
		tokens, _ := tokenizer.Tokenize()
		parser := NewParser(tokens, false)
		expr, err := parser.Parse()
		es := (expr[0]).(*ExpressionStatement)
		ps.Assert().IsType(&LiteralExpr{}, es.E)
		ps.Equal(0, err)
	})
	ps.Run("test with nil", func() {
		tokenizer := scanner.NewTokenizer("nil;")
		tokens, _ := tokenizer.Tokenize()
		parser := NewParser(tokens, false)
		expr, err := parser.Parse()
		es := (expr[0]).(*ExpressionStatement)
		ps.Assert().IsType(&LiteralExpr{}, es.E)
		ps.Equal(0, err)
	})
	ps.Run("test with string", func() {
		tokenizer := scanner.NewTokenizer("\"hello world\";")
		tokens, _ := tokenizer.Tokenize()
		parser := NewParser(tokens, false)
		expr, err := parser.Parse()
		es := (expr[0]).(*ExpressionStatement)
		ps.Assert().IsType(&LiteralExpr{}, es.E)
		ps.Equal(0, err)
	})
	ps.Run("test with number", func() {
		tokenizer := scanner.NewTokenizer("42.37;")
		tokens, _ := tokenizer.Tokenize()
		parser := NewParser(tokens, false)
		expr, err := parser.Parse()
		es := (expr[0]).(*ExpressionStatement)
		ps.Assert().IsType(&LiteralExpr{}, es.E)
		ps.Equal(0, err)
	})
	ps.Run("test with parenthesized expression", func() {
		tokenizer := scanner.NewTokenizer("(true);")
		tokens, _ := tokenizer.Tokenize()
		parser := NewParser(tokens, false)
		expr, err := parser.Parse()
		es := (expr[0]).(*ExpressionStatement)
		ps.Assert().IsType(&GroupingExpr{}, es.E)
		ps.Equal(0, err)
	})
	ps.Run("test with unary expression", func() {
		tokenizer := scanner.NewTokenizer("!true;")
		tokens, _ := tokenizer.Tokenize()
		parser := NewParser(tokens, false)
		expr, err := parser.Parse()
		es := (expr[0]).(*ExpressionStatement)
		ps.Assert().IsType(&UnaryExpr{}, es.E)
		ps.Equal(0, err)
	})
	ps.Run("test with binary expression", func() {
		tokenizer := scanner.NewTokenizer("1+2;")
		tokens, _ := tokenizer.Tokenize()
		parser := NewParser(tokens, false)
		expr, err := parser.Parse()
		es := (expr[0]).(*ExpressionStatement)
		ps.Assert().IsType(&BinaryExpr{}, es.E)
		ps.Equal(0, err)
	})
	ps.Run("test with comparison", func() {
		tokenizer := scanner.NewTokenizer("1<2;")
		tokens, _ := tokenizer.Tokenize()
		parser := NewParser(tokens, false)
		expr, err := parser.Parse()
		es := (expr[0]).(*ExpressionStatement)
		ps.Assert().IsType(&BinaryExpr{}, es.E)
		ps.Equal(0, err)
	})
	ps.Run("test with equality", func() {
		tokenizer := scanner.NewTokenizer("1==2;")
		tokens, _ := tokenizer.Tokenize()
		parser := NewParser(tokens, false)
		expr, err := parser.Parse()
		es := (expr[0]).(*ExpressionStatement)
		ps.Assert().IsType(&BinaryExpr{}, es.E)
		ps.Equal(0, err)
	})
	ps.Run("test with print", func() {
		tokenizer := scanner.NewTokenizer("print 1;")
		tokens, _ := tokenizer.Tokenize()
		parser := NewParser(tokens, false)
		expr, err := parser.Parse()
		ps.Assert().IsType(&PrintStatement{}, expr[0])
		ps.Equal(0, err)
	})
	ps.Run("parse with error", func() {
		tokenizer := scanner.NewTokenizer("(foo")
		tokens, _ := tokenizer.Tokenize()
		parser := NewParser(tokens, true)
		expr, err := parser.Parse()
		ps.Assert().IsType(&ExpressionStatement{}, expr[0])
		ps.Equal(65, err)
	})
	ps.Run("parse var declaration", func() {
		tokenizer := scanner.NewTokenizer("var a = 1;")
		tokens, _ := tokenizer.Tokenize()
		parser := NewParser(tokens, false)
		expr, err := parser.Parse()
		ps.Assert().IsType(&VarStatement{}, expr[0])
		ps.Equal(0, err)
	})
	ps.Run("parse assignment", func() {
		tokenizer := scanner.NewTokenizer("a = 1;")
		tokens, _ := tokenizer.Tokenize()
		parser := NewParser(tokens, false)
		expr, err := parser.Parse()
		ps.Assert().IsType(&ExpressionStatement{}, expr[0])
		ps.Equal(0, err)
	})
	ps.Run("parse blocks", func() {
		tokenizer := scanner.NewTokenizer("{print 1;var a = \"hello\";}")
		tokens, _ := tokenizer.Tokenize()
		parser := NewParser(tokens, false)
		expr, err := parser.Parse()
		ps.Assert().IsType(&BlockStatement{}, expr[0])
		ps.Equal(0, err)
	})
	ps.Run("parse if statement", func() {
		tokenizer := scanner.NewTokenizer("if (true) {print 1;}")
		tokens, _ := tokenizer.Tokenize()
		parser := NewParser(tokens, false)
		expr, err := parser.Parse()
		ps.Assert().IsType(&IfStatement{}, expr[0])
		ps.Equal(0, err)
	})
}
