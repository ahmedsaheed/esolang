package object

import (
	"fmt"
)

func newError(format string, a ...interface{}) Object {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func NewError(format string, a ...interface{}) Object {
	return newError(format, a...)
}

func newErrorFromTypings(err string) Object {
	return &Error{Message: err}
}

func arrayInvokables(method string, arr *Array, args ...Object) Object {
	name := "Array." + method
	if method == "count" || method == "length" {
		if err := CheckTypings(
			name, args,
			ExactArgsLength(0),
		); err != nil {
			return newErrorFromTypings(err.Error())
		}
		return &Integer{Value: int64(len(arr.Elements))}
	}

	if method == "sort" {
		if err := CheckTypings(
			name, args,
			ExactArgsLength(0),
		); err != nil {
			return newErrorFromTypings(err.Error())
		}
		return array_sort(arr)
	}

	if method == "get" {
		if err := CheckTypings(
			name, args,
			ExactArgsLength(1),
			WithTypes(INTEGER_OBJ),
		); err != nil {
			return newErrorFromTypings(err.Error())
		}

		index := args[0].(*Integer).Value
		max := int64(len(arr.Elements) - 1)

		if index < 0 || index > max {
			return &Null{}
		}

		return arr.Elements[index]
	}

	if method == "get_first" {
		if err := CheckTypings(
			name, args,
			ExactArgsLength(1),
			WithTypes(INTEGER_OBJ),
		); err != nil {
			return newErrorFromTypings(err.Error())
		}
		if len(arr.Elements) > 0 {
			return arr.Elements[0]
		}
		return &Null{}
	}

	if method == "get_last" {
		if err := CheckTypings(
			name, args,
			ExactArgsLength(1),
			WithTypes(INTEGER_OBJ),
		); err != nil {
			return newErrorFromTypings(err.Error())
		}
		if len(arr.Elements) > 0 {
			return arr.Elements[len(arr.Elements)-1]
		}
		return &Null{}
	}

	if method == "append" {
		if err := CheckTypings(
			name, args,
			ExactArgsLength(1),
			WithTypes(
				INTEGER_OBJ,
				STRING_OBJ,
				BOOLEAN_OBJ,
				NULL_OBJ,
				ARRAY_OBJ,
				HASH_OBJ,
				FUNCTION_OBJ,
			),
		); err != nil {
			return newErrorFromTypings(err.Error())
		}
		newArr := append(arr.Elements, args[0])
		return &Array{Elements: newArr}
	}
	if method == "fill" {
		/*
			fill(value)
			fill(value, start)
			fill(value, start, end)
		*/
		if err := CheckTypings(
			name, args,
			RangeOfArgs(1, 3),
		); err != nil {
			return newErrorFromTypings(err.Error())
		}
		return array_fill(arr, args...)
	}

	if method == "index_of" {
		if err := CheckTypings(
			name, args,
			ExactArgsLength(1),
		); err != nil {
			return newErrorFromTypings(err.Error())
		}

		target := args[0]

		for i, el := range arr.Elements {
			if el.Inspect() == target.Inspect() {
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
