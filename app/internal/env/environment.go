package env

type Env struct {
	values map[string]any
}

func NewEnv() *Env {
	return &Env{
		values: make(map[string]any),
	}
}

func (e *Env) Get(key string) any {
	return e.values[key]
}

func (e *Env) Put(key string, value any) {
	e.values[key] = value
}
