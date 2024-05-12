package object

import (
	"bytes"
	"esolang/lang-esolang/ast"
	"fmt"
	"hash/fnv"
	"strings"
)

type ObjectType string

const (
	STRING_OBJ       = "STRING"
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	FUNCTION_OBJ     = "FUNCTION"
	ERROR_OBJ        = "ERROR"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
	MODULE_TYPE      = "MODULE"
)

type Object interface {
	Type() ObjectType
	Inspect() string
	InvokeMethod(method string, env Environment, args ...Object) Object
}

// String wraps a single value to a string.
type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

func (s *String) InvokeMethod(method string, env Environment, args ...Object) Object {
	return stringInvokables(method, s, args...)
}

// Integer wraps a single value to an integer64.
type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) InvokeMethod(method string, env Environment, args ...Object) Object {
	return intInvokables(method, i)
}

// Boolean wraps a single value to a boolean.
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) InvokeMethod(method string, env Environment, args ...Object) Object {
	return nil
}

// Null represents the absence of a value.
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }
func (n *Null) InvokeMethod(method string, env Environment, args ...Object) Object {
	return nil
}

// ReturnValue wraps a single value to a return value.
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
func (rv *ReturnValue) InvokeMethod(method string, env Environment, args ...Object) Object {
	return nil
}

// Function wraps a block statement to a function.
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()

}
func (f *Function) InvokeMethod(method string, env Environment, args ...Object) Object {
	return nil
}

// Array wraps a list of objects to an array.
type Array struct {
	Elements []Object
}

func (ao *Array) Type() ObjectType { return ARRAY_OBJ }
func (ao *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, element := range ao.Elements {
		elements = append(elements, element.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
func (ao *Array) InvokeMethod(method string, env Environment, args ...Object) Object {
	return arrayInvokables(method, ao, args...)
}

// HashKey is a key for a hash.
type HashKey struct {
	Type  ObjectType
	Value uint64
}

// HashKey returns a hash key for a boolean.
func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

// HashKey returns a hash key for an integer.
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }

func (h *Hash) Inspect() string {
	var output bytes.Buffer

	pairs := []string{}

	for _, pair := range h.Pairs {
		// wrap string key & value with double quotes

		if pair.Key.Type() == STRING_OBJ && pair.Value.Type() == STRING_OBJ {
			pairs = append(pairs, fmt.Sprintf("\"%s\": \"%s\"", pair.Key.Inspect(), pair.Value.Inspect()))
			continue
		} else if pair.Value.Type() == STRING_OBJ {
			pairs = append(pairs, fmt.Sprintf("%s: \"%s\"", pair.Key.Inspect(), pair.Value.Inspect()))
			continue
		} else if pair.Key.Type() == STRING_OBJ {
			pairs = append(pairs, fmt.Sprintf("\"%s\": %s", pair.Key.Inspect(), pair.Value.Inspect()))
			continue
		}

		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	output.WriteString("{")
	output.WriteString(strings.Join(pairs, ", "))
	output.WriteString("}")

	return output.String()
}
func (h *Hash) InvokeMethod(method string, env Environment, args ...Object) Object {
	return hashInvokables(method, h)
}

type Hashable interface {
	HashKey() HashKey
}

// Error wraps a single value to an error.
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }
func (e *Error) InvokeMethod(method string, env Environment, args ...Object) Object {
	return nil
}

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }
func (b *Builtin) InvokeMethod(method string, env Environment, args ...Object) Object {
	return nil
}

type Hasher interface {
	Hash() HashKey
}

type Module struct {
	Name  string
	Attrs Object
}

func (m Module) Bool() bool {
	return true
}

func (m Module) Get(index Object) (Object, error) {
	key, ok := index.(Hasher)
	if !ok {
		return nil, fmt.Errorf("invalid module attribute %s", index.Type())
	}

	attr, found := m.Attrs.(*Hash).Pairs[key.Hash()]
	if !found {
		return &Null{}, nil
	}

	return attr.Value, nil
}

func (m Module) Compare(other Object) int {
	return 1
}

func (m Module) String() string {
	return m.Inspect()
}

// Type returns the type of the object
func (m Module) Type() ObjectType { return MODULE_TYPE }

// Inspect returns a stringified version of the object for debugging
func (m Module) Inspect() string { return fmt.Sprintf("<module '%s'>", m.Name) }

// InvokeMethod implements Object.
func (m *Module) InvokeMethod(method string, env Environment, args ...Object) Object {
	panic("unimplemented")
}
