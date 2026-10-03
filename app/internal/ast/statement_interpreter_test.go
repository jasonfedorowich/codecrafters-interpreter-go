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
	interpreter *Interpreter
	ctrl        *gomock.Controller
	p           *mocks.MockPrinter
}

func (is *StatementInterpreterTestSuite) SetupSubTest() {
	is.ctrl = gomock.NewController(is.T())
	is.p = mocks.NewMockPrinter(is.ctrl)
	is.interpreter = &Interpreter{
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
		_ = is.interpreter.InterpretStatements(statements)
	})
	is.Run("test print with number", func() {
		is.p.EXPECT().Print("124").Times(1)
		tokenizer := scanner.NewTokenizer("print (49 * 2 + 75 * 2) / (2);")
		tokens, _ := tokenizer.Tokenize()
		parser := parse.NewParser(tokens, true)
		statements, _ := parser.Parse()
		_ = is.interpreter.InterpretStatements(statements)
	})
	is.Run("test invalid print", func() {
		tokenizer := scanner.NewTokenizer("print false * (24 + 68);")
		tokens, _ := tokenizer.Tokenize()
		parser := parse.NewParser(tokens, true)
		statements, _ := parser.Parse()
		_ = is.interpreter.InterpretStatements(statements)
	})
	is.Run("test var declaration", func() {
		is.p.EXPECT().Print("1").Times(1)
		tokenizer := scanner.NewTokenizer("var a = 1; print a;")
		tokens, _ := tokenizer.Tokenize()
		parser := parse.NewParser(tokens, false)
		statements, _ := parser.Parse()
		_ = is.interpreter.InterpretStatements(statements)
	})
	is.Run("test statements and declarations", func() {
		is.p.EXPECT().Print("152").Times(1)
		is.p.EXPECT().Print("228").Times(1)
		tokenizer := scanner.NewTokenizer("var quz = 76;\n\t var hello = 76;\n\t print quz + hello;\n\t var world = 76;\n\t\tprint quz + hello + world;")
		tokens, _ := tokenizer.Tokenize()
		parser := parse.NewParser(tokens, false)
		statements, _ := parser.Parse()
		_ = is.interpreter.InterpretStatements(statements)
	})
	is.Run("test undefined variable", func() {
		tokenizer := scanner.NewTokenizer("print quz;")
		tokens, _ := tokenizer.Tokenize()
		parser := parse.NewParser(tokens, false)
		statements, _ := parser.Parse()
		_ = is.interpreter.InterpretStatements(statements)
	})
	is.Run("test undefined variable", func() {
		is.p.EXPECT().Print("nil").Times(1)
		tokenizer := scanner.NewTokenizer("var foo;\nprint foo;")
		tokens, _ := tokenizer.Tokenize()
		parser := parse.NewParser(tokens, false)
		statements, _ := parser.Parse()
		_ = is.interpreter.InterpretStatements(statements)

	})
	is.Run("test assignment", func() {
		is.p.EXPECT().Print("1").Times(1)
		tokenizer := scanner.NewTokenizer("var foo;\nfoo = 1;\nprint foo;")
		tokens, _ := tokenizer.Tokenize()
		parser := parse.NewParser(tokens, false)
		statements, _ := parser.Parse()
		_ = is.interpreter.InterpretStatements(statements)

	})
	is.Run("test block statements", func() {
		is.p.EXPECT().Print("before").Times(1)
		is.p.EXPECT().Print("after").Times(1)
		tokenizer := scanner.NewTokenizer("{\n    var foo = \"before\";\n    print foo;\n}\n{\n    var foo = \"after\";\n    print foo;\n}")
		tokens, _ := tokenizer.Tokenize()
		parser := parse.NewParser(tokens, false)
		statements, _ := parser.Parse()
		_ = is.interpreter.InterpretStatements(statements)
	})

}
