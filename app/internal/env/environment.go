package env

import "fmt"

type Env struct {
	values map[string]any
	e      *Env
}

func NewEnv() *Env {
	return &Env{
		values: make(map[string]any),
	}
}

func NewWrappedEnv(e *Env) *Env {
	return &Env{
		values: make(map[string]any),
		e:      e,
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
	if e.e != nil {
		return e.e.Assign(key, value)
	}
	return fmt.Errorf("key does not exist: %s", key)
}

func (e *Env) CheckAndGet(key string) (v any, ok bool) {
	v, ok = e.values[key]
	return
}

func (e *Env) Get(key string) any {
	if e.Contains(key) {
		return e.values[key]
	}
	if e.e != nil {
		return e.e.Get(key)
	}
	return nil
}

func (e *Env) Put(key string, value any) {
	e.values[key] = value
}
