package ast

import (
	"testing"

	"github.com/codecrafters-io/interpreter-starter-go/app/internal/ast/mocks"
	"github.com/codecrafters-io/interpreter-starter-go/app/internal/env"
	"github.com/codecrafters-io/interpreter-starter-go/app/internal/error_code"
	"github.com/codecrafters-io/interpreter-starter-go/app/internal/parse"
	"github.com/codecrafters-io/interpreter-starter-go/app/internal/scanner"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type StatementInterpreterTestSuite struct {
	suite.Suite
	interpreter *StatementInterpreter
	ctrl        *gomock.Controller
	p           *mocks.MockPrinter
}

func (is *StatementInterpreterTestSuite) SetupSubTest() {
	is.ctrl = gomock.NewController(is.T())
	is.p = mocks.NewMockPrinter(is.ctrl)
	is.interpreter = &StatementInterpreter{
		eh: &error_code.ErrorHandler{},
		p:  is.p,
		e:  env.NewEnv(),
	}
}

func (is *StatementInterpreterTestSuite) TearDownTest() {

}

func TestStatementInterpreterTestSuite(t *testing.T) {
	suite.Run(t, new(StatementInterpreterTestSuite))
}

func (is *StatementInterpreterTestSuite) TestStatementInterpreter() {
	is.Run("test print statement", func() {
		is.p.EXPECT().Print("-7").Times(1)
		tokenizer := scanner.NewTokenizer("print -10*10+(100-7);")
		tokens, _ := tokenizer.Tokenize()
		parser := parse.NewParser(tokens, true)
		statements, _ := parser.Parse()
		_ = is.interpreter.Interpret(statements)
	})
	is.Run("test print with number", func() {
		is.p.EXPECT().Print("124").Times(1)
		tokenizer := scanner.NewTokenizer("print (49 * 2 + 75 * 2) / (2);")
		tokens, _ := tokenizer.Tokenize()
		parser := parse.NewParser(tokens, true)
		statements, _ := parser.Parse()
		_ = is.interpreter.Interpret(statements)
	})
	is.Run("test invalid print", func() {
		tokenizer := scanner.NewTokenizer("print false * (24 + 68);")
		tokens, _ := tokenizer.Tokenize()
		parser := parse.NewParser(tokens, true)
		statements, _ := parser.Parse()
		_ = is.interpreter.Interpret(statements)
	})
	is.Run("test var declaration", func() {
		is.p.EXPECT().Print("1").Times(1)
		tokenizer := scanner.NewTokenizer("var a = 1; print a;")
		tokens, _ := tokenizer.Tokenize()
		parser := parse.NewParser(tokens, false)
		statements, _ := parser.Parse()
		_ = is.interpreter.Interpret(statements)
	})

}
