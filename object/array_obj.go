package object

import (
	"fmt"
	"sort"
)

func newError(format string, a ...interface{}) Object {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func arrayInvokables(method string, arr *Array, args ...Object) Object {
	if method == "count" || method == "length" {
		return &Integer{Value: int64(len(arr.Elements))}
	}

	if method == "sort" {
		elements := arr.Elements
		intsort, intArr, strsort, strArr := _checkSortType(elements)

		if !intsort && !strsort {
			return newError("sort method only works on arrays of integers or strings")
		}

		if intsort {
			sort.Ints(intArr)
			elements = make([]Object, len(intArr))
			for i, v := range intArr {
				elements[i] = &Integer{Value: int64(v)}
			}
			return &Array{Elements: elements}
		}
		if strsort {
			sort.Strings(strArr)
			elements = make([]Object, len(strArr))
			for i, v := range strArr {
				elements[i] = &String{Value: v}
			}
			return &Array{Elements: elements}
		}
	}

	if method == "get" {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1", len(args))
		}

		if args[0].Type() != INTEGER_OBJ {
			return newError("argument to `get` must be INTEGER, got %s", args[0].Type())
		}

		index := args[0].(*Integer).Value
		max := int64(len(arr.Elements) - 1)

		if index < 0 || index > max {
			return &Null{}
		}

		return arr.Elements[index]
	}

	if method == "get_first" {
		if len(arr.Elements) > 0 {
			return arr.Elements[0]
		}
		return &Null{}
	}

	if method == "get_last" {
		if len(arr.Elements) > 0 {
			return arr.Elements[len(arr.Elements)-1]
		}
		return &Null{}
	}

	if method == "append" {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1", len(args))
		}
		newArr := append(arr.Elements, args[0])
		return &Array{Elements: newArr}
	}

	if method == "pop" {

	}

	if method == "fill" {
		/*
			fill(value)
			fill(value, start)
			fill(value, start, end)
		*/

		if len(args) < 1 || len(args) > 3 {
			return newError("wrong number of arguments. got=%d, want=1, 2 or 3", len(args))
		}

		return array_fill(method, arr, args...)

	}

	if method == "index_of" {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1", len(args))
		}

		target := args[0]

		for i, el := range arr.Elements {
			if el == target {
				return &Integer{Value: int64(i)}
			}
		}

		return &Null{}
	}

	return nil
}

func hashInvokables(method string, h *Hash) Object {
	if method == "count" || method == "length" {
		return &Integer{Value: int64(len(h.Pairs))}
	}

	if method == "keys" {
		keys := &Array{}
		for _, pair := range h.Pairs {
			keys.Elements = append(keys.Elements, pair.Key)
		}
		return keys
	}

	if method == "values" {
		values := &Array{}
		for _, pair := range h.Pairs {
			values.Elements = append(values.Elements, pair.Value)
		}
		return values
	}

	if method == "entries" {
		entries := &Array{}
		for _, pair := range h.Pairs {
			entry := &Array{Elements: []Object{pair.Key, pair.Value}}
			entries.Elements = append(entries.Elements, entry)
		}
		return entries
	}

	if method == "to_string" {
		return &String{Value: h.Inspect()}
	}

	return nil
}
