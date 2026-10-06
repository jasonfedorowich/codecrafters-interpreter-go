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
		e:  env.NewEnv(loadNativeFunctions()),
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
	is.Run("test scopes", func() {
		is.p.EXPECT().Print(gomock.Any()).AnyTimes()
		tokenizer := scanner.NewTokenizer("var bar = \"global bar\";\nvar world = \"global world\";\nvar hello = \"global hello\";\n{\n  var bar = \"outer bar\";\n  var world = \"outer world\";\n  {\n    var bar = \"inner bar\";\n    print bar;\n    print world;\n    print hello;\n  }\n  print bar;\n  print world;\n  print hello;\n}\nprint bar;\nprint world;\nprint hello;")
		tokens, _ := tokenizer.Tokenize()
		parser := parse.NewParser(tokens, false)
		statements, _ := parser.Parse()
		err := is.interpreter.InterpretStatements(statements)
		is.Assert().Equal(0, err)
	})
	is.Run("test if statement", func() {
		is.p.EXPECT().Print("1").Times(1)
		tokenizer := scanner.NewTokenizer("if (true) {print 1;}")
		tokens, _ := tokenizer.Tokenize()
		parser := parse.NewParser(tokens, false)
		statements, _ := parser.Parse()
		err := is.interpreter.InterpretStatements(statements)
		is.Assert().Equal(0, err)
	})
	is.Run("logical OR", func() {
		tokenizer := scanner.NewTokenizer("true or false;")
		tokens, _ := tokenizer.Tokenize()
		parser := parse.NewParser(tokens, true)
		statements, _ := parser.Parse()
		res, err := is.interpreter.InterpretExpr(statements[0])
		is.Equal("true", res)
		is.Equal(0, err)
	})
	is.Run("while statement", func() {
		is.p.EXPECT().Print(gomock.Any()).AnyTimes()
		tokenizer := scanner.NewTokenizer("var foo = 0;\nwhile (foo < 3) print foo = foo + 1;")
		tokens, _ := tokenizer.Tokenize()
		parser := parse.NewParser(tokens, false)
		statements, _ := parser.Parse()
		err := is.interpreter.InterpretStatements(statements)
		is.Assert().Equal(0, err)
	})

	is.Run("for statement", func() {
		is.p.EXPECT().Print(gomock.Any()).AnyTimes()
		tokenizer := scanner.NewTokenizer("for (var i = 0; i < 3; i = i + 1) print i;")
		tokens, _ := tokenizer.Tokenize()
		parser := parse.NewParser(tokens, false)
		statements, _ := parser.Parse()
		err := is.interpreter.InterpretStatements(statements)
		is.Assert().Equal(0, err)
	})
	is.Run("native function call", func() {
		is.p.EXPECT().Print(gomock.Any()).AnyTimes()
		tokenizer := scanner.NewTokenizer("print clock();")
		tokens, _ := tokenizer.Tokenize()
		parser := parse.NewParser(tokens, false)
		statements, _ := parser.Parse()
		err := is.interpreter.InterpretStatements(statements)
		is.Assert().Equal(0, err)
	})

	is.Run("function define and call", func() {
		is.p.EXPECT().Print(gomock.Any()).AnyTimes()
		tokenizer := scanner.NewTokenizer("fun foo() { print 1; }\n foo();")
		tokens, _ := tokenizer.Tokenize()
		parser := parse.NewParser(tokens, false)
		statements, _ := parser.Parse()
		err := is.interpreter.InterpretStatements(statements)
		is.Assert().Equal(0, err)
	})

	is.Run("function return", func() {
		is.p.EXPECT().Print(gomock.Any()).AnyTimes()
		tokenizer := scanner.NewTokenizer("fun foo() { return 1; }\n print foo();")
		tokens, _ := tokenizer.Tokenize()
		parser := parse.NewParser(tokens, false)
		statements, _ := parser.Parse()
		err := is.interpreter.InterpretStatements(statements)
		is.Assert().Equal(0, err)
	})
	is.Run("test with early return", func() {
		is.p.EXPECT().Print("nil").AnyTimes()
		tokenizer := scanner.NewTokenizer("fun foo() { return; print \"bad\"; }\n print foo();")
		tokens, _ := tokenizer.Tokenize()
		parser := parse.NewParser(tokens, false)
		statements, _ := parser.Parse()
		err := is.interpreter.InterpretStatements(statements)
		is.Assert().Equal(0, err)
	})

}
