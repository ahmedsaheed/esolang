package builtins

import (
	"esolang/lang-esolang/object"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var NULL = &object.Null{}

var Builtins = map[string]*object.Builtin{
	"Set": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) == 0 {
				return &object.Set{Elements: []object.Object{}}
			}
			return newError("wrong number of arguments. got=%d, want=0", len(args))
		},
	},
	"count": &object.Builtin{
		Fn: arrayLen,
	},
	"array_getFirst": &object.Builtin{
		Fn: arrayGetFirst,
	},
	"array_getLast": &object.Builtin{
		Fn: arrayGetLast,
	},
	"array_get": &object.Builtin{
		Fn: arrayGet,
	},
	"array_pop": &object.Builtin{
		Fn: arrayPop,
	},
	"array_append": &object.Builtin{
		Fn: arrayAppend,
	},
	"array_new": &object.Builtin{
		Fn: arrayNew,
	},

	"array_indexOf": &object.Builtin{
		Fn: arrayIndexOf,
	},
	"array_sort": &object.Builtin{
		Fn: arraySort,
	},
	"array_rest": &object.Builtin{
		Fn: arrayRest,
	},
	"println": &object.Builtin{
		Fn: Println,
	},
	"print": &object.Builtin{
		Fn: Print,
	},
	"ReadFile": &object.Builtin{
		Fn: readFile,
	},
	"WriteFile": &object.Builtin{
		Fn: writeFile,
	},
	"Http": &object.Builtin{
		Fn: _http,
	},
	"type_of": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			return &object.String{Value: string(args[0].Type())}
		},
	},
	"_upperCase": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError("argument to `upperCase` must be STRING, got %s", args[0].Type())
			}
			upperCased := strings.ToUpper(args[0].(*object.String).Value)
			return &object.String{Value: upperCased}
		},
	},

	"_lowerCase": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError("argument to `lowerCase` must be STRING, got %s", args[0].Type())
			}
			lowerCased := strings.ToLower(args[0].(*object.String).Value)
			return &object.String{Value: lowerCased}
		},
	},

	"math_randomInt": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			rand.Seed(time.Now().UnixNano())
			return &object.Integer{Value: rand.Int63()}
		},
	},

	"math_randomFloat": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			rand.Seed(time.Now().UnixNano())
			return &object.Float{Value: rand.Float64()}
		},
	},
	"_contains": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.STRING_OBJ || args[1].Type() != object.STRING_OBJ {
				return newError("arguments to `contains` must be STRING, got %s", args[0].Type())
			}
			contains := strings.Contains(args[0].(*object.String).Value, args[1].(*object.String).Value)
			return &object.Boolean{Value: contains}
		},
	},
}

func RegisterBuiltin(name string, fun object.BuiltinFunction) {
	Builtins[name] = &object.Builtin{Fn: fun}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
