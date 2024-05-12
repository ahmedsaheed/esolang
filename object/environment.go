package object

import "unicode"

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func (e *Environment) ExportedHash() *Hash {
	pairs := make(map[HashKey]HashPair)
	for k, v := range e.store {
		if unicode.IsUpper(rune(k[0])) {
			s := String{Value: k}
			pairs[s.HashKey()] = HashPair{Key: &s, Value: v}
		}
	}
	return &Hash{Pairs: pairs}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}
