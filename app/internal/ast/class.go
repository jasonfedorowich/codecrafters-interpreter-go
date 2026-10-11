package ast

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/app/internal/scanner"
)

type class struct {
	name string
}

type instance struct {
	cl     *class
	fields map[string]any
}

func newClass(name string) *class {
	return &class{
		name: name,
	}
}

func newInstance(cl *class) *instance {
	return &instance{
		cl:     cl,
		fields: make(map[string]any),
	}
}

func (c *class) call(inter *Interpreter, args []any) any {
	ins := newInstance(c)
	return ins
}

func (c *class) arity() int {
	return 0
}

func (c *class) String() string {
	return c.name
}

func (i *instance) String() string {
	return i.cl.name + " instance"
}

func (i *instance) get(name *scanner.Token) (any, error) {
	if v, ok := i.fields[name.Lexeme]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("Undefined property ' + %s + '.", name)
}

func (i *instance) set(name *scanner.Token, value any) {
	i.fields[name.Lexeme] = value
}
