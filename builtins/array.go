package builtins

import (
	"esolang/lang-esolang/object"
	"sort"
)

func arrayLen(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	switch arg := args[0].(type) {
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Elements))}
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	default:
		return newError("argument to `count` not supported, got %s", args[0].Type())
	}
}

func arrayGetFirst(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `getFirst` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	if len(arr.Elements) > 0 {
		return arr.Elements[0]
	}
	return NULL
}

func arrayGetLast(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `getLast` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)
	if length > 0 {
		return arr.Elements[length-1]
	}
	return NULL
}

func arrayGet(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2",
			len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("first argument to `get` must be ARRAY, got %s", args[0].Type())
	}

	if args[1].Type() != object.INTEGER_OBJ {
		return newError("second argument to `get` must be INTEGER, got %s", args[1].Type())
	}

	arr := args[0].(*object.Array)
	index := args[1].(*object.Integer).Value
	length := int64(len(arr.Elements))

	if index < 0 || index >= length {
		return NULL
	}

	return arr.Elements[index]
}

func arrayIndexOf(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2",
			len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("first argument to `find` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	target := args[1]

	for i, el := range arr.Elements {
		if el.Inspect() == target.Inspect() {
			return &object.Integer{Value: int64(i)}
		}
	}
	return NULL
}

func arrayPop(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `pop` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)
	if length > 0 {
		last := arr.Elements[length-1]
		arr.Elements = arr.Elements[:length-1]
		return last
	}
	return arr
}

func arrayAppend(args ...object.Object) object.Object {
	if len(args) < 2 {
		return newError("wrong number of arguments. got=%d, want=at least 2",
			len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("first argument to `push` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	newElements := args[1:]

	arr.Elements = append(arr.Elements, newElements...)
	return arr
}

func arrayNew(args ...object.Object) object.Object {
	return &object.Array{Elements: args}
}

func arraySort(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("first argument to `sort` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	elements := arr.Elements

	var intSort bool
	var strSort bool

	// check all elements are integers
	for _, el := range elements {
		if el.Type() == object.INTEGER_OBJ {
			intSort = true
		} else {
			intSort = false
			break
		}
	}

	// check all elements are strings
	for _, el := range elements {
		if el.Type() == object.STRING_OBJ {
			strSort = true
		} else {
			strSort = false
			break
		}
	}

	if !intSort && !strSort {
		return newError("Sorting only supports primitive types (integers and strings)")
	}

	if intSort {
		sort.Slice(elements, func(i, j int) bool {
			return elements[i].(*object.Integer).Value < elements[j].(*object.Integer).Value
		})
	}

	if strSort {
		sort.Slice(elements, func(i, j int) bool {
			return elements[i].(*object.String).Value < elements[j].(*object.String).Value
		})
	}

	return NULL
}

func arrayRest(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("first argument to `rest` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)

	if length > 0 {
		newElements := make([]object.Object, length-1, length-1)
		copy(newElements, arr.Elements[1:length])
		return &object.Array{Elements: newElements}
	}

	return NULL
}
