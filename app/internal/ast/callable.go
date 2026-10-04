package ast

import "time"

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
