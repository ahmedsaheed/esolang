package builtins

import (
	"esolang/lang-esolang/object"

	"fmt"
)

var NULL = &object.Null{}

var Builtins = map[string]*object.Builtin{
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
	"HttpServe": &object.Builtin{
		Fn: _httpServeAndListen,
	},
	"typeOf": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			return &object.String{Value: string(args[0].Type())}
		},
	},
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
