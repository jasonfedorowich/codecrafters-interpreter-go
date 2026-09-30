package ast

import (
	"testing"

	"github.com/codecrafters-io/interpreter-starter-go/app/internal/ast/mocks"
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

func (is *StatementInterpreterTestSuite) SetupTest() {
	is.ctrl = gomock.NewController(is.T())
	is.p = mocks.NewMockPrinter(is.ctrl)
	is.interpreter = &StatementInterpreter{
		eh: &error_code.ErrorHandler{},
		p:  is.p,
	}
}

func (is *StatementInterpreterTestSuite) TearDownTest() {

}

func TestStatementInterpreterTestSuite(t *testing.T) {
	suite.Run(t, new(StatementInterpreterTestSuite))
}

func (is *StatementInterpreterTestSuite) TestStatementInterpreter() {
	is.Run("test print statement", func() {
		is.p.EXPECT().Print("-7.0").Times(1)
		tokenizer := scanner.NewTokenizer("print -10*10+(100-7);")
		tokens, _ := tokenizer.Tokenize()
		parser := parse.NewParser(tokens, true)
		statements, _ := parser.Parse()
		_ = is.interpreter.Interpret(statements)
	})

}
