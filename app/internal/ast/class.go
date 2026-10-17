package ast

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/app/internal/scanner"
)

type class struct {
	name       string
	methods    map[string]*function
	superClass *class
}

func newClass(name string, methods map[string]*function, superClass *class) *class {
	return &class{
		name:       name,
		methods:    methods,
		superClass: superClass,
	}
}

func (c *class) call(inter *Interpreter, args []any) any {
	ins := newInstance(c)
	initializer := c.findMethod("init")
	if initializer != nil {
		initializer.bind(ins).call(inter, args)
	}
	return ins
}

func (c *class) arity() int {
	initializer := c.findMethod("init")
	if initializer == nil {
		return 0
	}
	return initializer.arity()
}

func (c *class) String() string {
	return c.name
}

func (c *class) findMethod(lexeme string) *function {
	if c.methods != nil && c.methods[lexeme] != nil {
		return c.methods[lexeme]
	}
	if c.superClass != nil {
		return c.superClass.findMethod(lexeme)
	}
	return nil
}

type instance struct {
	cl     *class
	fields map[string]any
}

func newInstance(cl *class) *instance {
	return &instance{
		cl:     cl,
		fields: make(map[string]any),
	}
}

func (i *instance) String() string {
	return i.cl.name + " instance"
}

func (i *instance) get(name *scanner.Token) (any, error) {
	if v, ok := i.fields[name.Lexeme]; ok {
		return v, nil
	}
	method := i.cl.findMethod(name.Lexeme)
	if method != nil {
		return method.bind(i), nil
	}
	return nil, fmt.Errorf("Undefined property ' + %s + '.", name)
}

func (i *instance) set(name *scanner.Token, value any) {
	i.fields[name.Lexeme] = value
}
