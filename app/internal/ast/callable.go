package ast

import (
	"fmt"
	"time"

	"github.com/codecrafters-io/interpreter-starter-go/app/internal/env"
	"github.com/codecrafters-io/interpreter-starter-go/app/internal/parse"
)

type returnValue struct {
	v any
}

func (r *returnValue) Error() string {
	return fmt.Sprintf("return: %v", r.v)
}

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
	closure     *env.Env
	inInit      bool
}

func newFunction(f *parse.FuncStatement, e *env.Env, inInit bool) *function {
	return &function{statementFn: f,
		closure: e}
}

func (f *function) call(interpreter *Interpreter, args []any) (r any) {
	defer func() {
		if e := recover(); e != nil {
			if v, ok := e.(*returnValue); ok {
				if f.inInit {
					r = f.closure.GetAt(0, "this")
					return
				}
				r = v.v
				return
			}
			panic(e)
		}
	}()
	e := env.NewWrappedEnv(f.closure)
	for i, a := range args {
		e.Put(f.statementFn.Params[i].Lexeme, a)
	}
	interpreter.blockStatement(f.statementFn.Body, e)
	if !f.inInit {
		return f.closure.GetAt(0, "this")
	}
	return
}

func (f *function) arity() int {
	return len(f.statementFn.Params)
}

func (f *function) String() string {
	return fmt.Sprintf("<fn %s>", f.statementFn.Name.Lexeme)
}

func (f *function) bind(i *instance) *function {
	e := env.NewWrappedEnv(f.closure)
	e.Put("this", i)
	return newFunction(f.statementFn, e, f.inInit)
}
