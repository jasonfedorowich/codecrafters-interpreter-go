package ast

import (
	"testing"

	"github.com/codecrafters-io/interpreter-starter-go/app/internal/parse"
	"github.com/codecrafters-io/interpreter-starter-go/app/internal/scanner"
	"github.com/stretchr/testify/suite"
)

type InterpreterTestSuite struct {
	suite.Suite
	interpreter *Interpreter
}

func (is *InterpreterTestSuite) SetupTest() {
	is.interpreter = NewInterpreter()
}

func (is *InterpreterTestSuite) TearDownTest() {

}

func TestInterpreterSuite(t *testing.T) {
	suite.Run(t, new(InterpreterTestSuite))
}

func (is *InterpreterTestSuite) TestInterpreter() {
	is.Run("evaluate literal", func() {
		tokenizer := scanner.NewTokenizer("true")
		tokens, _ := tokenizer.Tokenize()
		parser := parse.NewParser(tokens)
		expr, _ := parser.Parse()
		res, _ := is.interpreter.Interpret(expr)
		is.Equal("true", res)
	})
	is.Run("evaluate string", func() {
		tokenizer := scanner.NewTokenizer("\"hello world!\"")
		tokens, _ := tokenizer.Tokenize()
		parser := parse.NewParser(tokens)
		expr, _ := parser.Parse()
		res, _ := is.interpreter.Interpret(expr)
		is.Equal("hello world!", res)
	})
	is.Run("evaluate double", func() {
		tokenizer := scanner.NewTokenizer("10.0")
		tokens, _ := tokenizer.Tokenize()
		parser := parse.NewParser(tokens)
		expr, _ := parser.Parse()
		res, _ := is.interpreter.Interpret(expr)
		is.Equal("10", res)
	})
	is.Run("evaluate group", func() {
		tokenizer := scanner.NewTokenizer("(10)")
		tokens, _ := tokenizer.Tokenize()
		parser := parse.NewParser(tokens)
		expr, _ := parser.Parse()
		res, _ := is.interpreter.Interpret(expr)
		is.Equal("10", res)
	})
	is.Run("evaluate unary", func() {
		tokenizer := scanner.NewTokenizer("(-10)")
		tokens, _ := tokenizer.Tokenize()
		parser := parse.NewParser(tokens)
		expr, _ := parser.Parse()
		res, _ := is.interpreter.Interpret(expr)
		is.Equal("-10", res)
	})
	is.Run("evaluate unary nil", func() {
		tokenizer := scanner.NewTokenizer("!nil")
		tokens, _ := tokenizer.Tokenize()
		parser := parse.NewParser(tokens)
		expr, _ := parser.Parse()
		res, _ := is.interpreter.Interpret(expr)
		is.Equal("true", res)
	})
	is.Run("strings equality", func() {
		tokenizer := scanner.NewTokenizer("\"hello\" == \"hello\"")
		tokens, _ := tokenizer.Tokenize()
		parser := parse.NewParser(tokens)
		expr, _ := parser.Parse()
		res, _ := is.interpreter.Interpret(expr)
		is.Equal("true", res)
	})
	is.Run("strings inequality", func() {
		tokenizer := scanner.NewTokenizer("\"baz\" != \"world\"")
		tokens, _ := tokenizer.Tokenize()
		parser := parse.NewParser(tokens)
		expr, _ := parser.Parse()
		res, _ := is.interpreter.Interpret(expr)
		is.Equal("true", res)
	})
	is.Run("unary ! errors", func() {
		tokenizer := scanner.NewTokenizer("(!!39)")
		tokens, _ := tokenizer.Tokenize()
		parser := parse.NewParser(tokens)
		expr, _ := parser.Parse()
		res, err := is.interpreter.Interpret(expr)
		is.Equal("false", res)
		is.Equal(0, err)
	})
}
