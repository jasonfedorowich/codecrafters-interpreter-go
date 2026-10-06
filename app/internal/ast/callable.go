package ast

import (
	"fmt"
	"time"

	"github.com/codecrafters-io/interpreter-starter-go/app/internal/env"
	"github.com/codecrafters-io/interpreter-starter-go/app/internal/parse"
)

type Callable interface {
	call(interpreter *Interpreter, args []any) any
	arity() int
	String() string
}

type clock struct{}

func (c *clock) String() string {
	return "<native fn>"
}

func (c *clock) arity() int {
	return 0
}

func (c *clock) call(interpreter *Interpreter, args []any) any {
	return float64(time.Now().UnixMilli()) / 1000.0
}

func loadNativeFunctions() map[string]any {
	return map[string]any{"clock": &clock{}}
}

type function struct {
	statementFn *parse.FuncStatement
}

func newFunction(f *parse.FuncStatement) Callable {
	return &function{statementFn: f}
}

func (f *function) call(interpreter *Interpreter, args []any) any {
	e := env.NewWrappedEnv(interpreter.e)
	for i, a := range args {
		e.Put(f.statementFn.Params[i].Lexeme, a)
	}
	return interpreter.blockStatement(f.statementFn.Body.Statements, e)
}

func (f *function) arity() int {
	return len(f.statementFn.Params)
}

func (f *function) String() string {
	return fmt.Sprintf("<fn %s>", f.statementFn.Name.Lexeme)
}
