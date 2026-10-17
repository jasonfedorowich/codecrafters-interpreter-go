package env

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/app/internal/scanner"
)

type Env struct {
	values map[string]any
	Outer  *Env
}

func NewEnv(nativeFunctions map[string]any) *Env {
	values := make(map[string]any)
	for key, value := range nativeFunctions {
		values[key] = value
	}
	return &Env{
		values: values,
	}
}

func NewWrappedEnv(e *Env) *Env {
	return &Env{
		values: make(map[string]any),
		Outer:  e,
	}
}

func (e *Env) Contains(key string) bool {
	_, ok := e.values[key]
	return ok
}

func (e *Env) Assign(key string, value any) error {
	if e.Contains(key) {
		e.values[key] = value
		return nil
	}
	if e.Outer != nil {
		return e.Outer.Assign(key, value)
	}
	return fmt.Errorf("key does not exist: %s", key)
}

func (e *Env) CheckAndGet(key string) (v any, ok bool) {
	v, ok = e.values[key]
	return
}

func (e *Env) Get(key string) (any, bool) {
	if e.Contains(key) {
		v, ok := e.values[key]
		return v, ok
	}
	if e.Outer != nil {
		return e.Outer.Get(key)
	}
	return nil, false
}

func (e *Env) Put(key string, value any) {
	e.values[key] = value
}

func (e *Env) GetAt(d int, name string) any {
	return e.ancestor(d).values[name]
}

func (e *Env) AssignAt(d int, name *scanner.Token, value any) {
	e.ancestor(d).values[name.Lexeme] = value
}

func (e *Env) ancestor(d int) *Env {
	env := e
	for i := 0; i < d; i++ {
		env = env.Outer
	}
	return env
}
