package ast

type class struct {
	name string
}

type instance struct {
	cl *class
}

func newClass(name string) *class {
	return &class{
		name: name,
	}
}

func newInstance(cl *class) *instance {
	return &instance{
		cl: cl,
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
