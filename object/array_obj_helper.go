package object

import (
	"sort"
)

func _checkSortType(elements []Object) (bool, []int, bool, []string) {
	var intSort bool
	var strSort bool
	var intArr []int
	var strArr []string

	for _, el := range elements {
		if el.Type() == INTEGER_OBJ {
			intSort = true
			intArr = append(intArr, int(el.(*Integer).Value))
		} else {
			intSort = false
			intArr = nil
			break
		}
	}

	for _, el := range elements {
		if el.Type() == STRING_OBJ {
			strSort = true
			strArr = append(strArr, el.(*String).Value)
		} else {
			strSort = false
			strArr = nil
			break
		}
	}

	return intSort, intArr, strSort, strArr
}

func array_fill(arr *Array, args ...Object) Object {
	/*
		fill(value)
		fill(value, start)
		fill(value, start, end)
	*/

	value := args[0]
	start := 0
	end := len(arr.Elements)

	if len(args) == 1 {
		for i := range arr.Elements {
			arr.Elements[i] = value
		}
		return arr
	}

	if len(args) == 2 {
		if args[1].Type() != INTEGER_OBJ {
			return NewError("argument to `start` must be INTEGER, got %s", args[1].Type())
		}
		start = int(args[1].(*Integer).Value)
		for i := start; i < len(arr.Elements); i++ {
			arr.Elements[i] = value
		}
		return arr
	}

	if len(args) == 3 {
		if args[1].Type() != INTEGER_OBJ {
			return NewError("argument to `start` must be INTEGER, got %s", args[1].Type())
		}
		start = int(args[1].(*Integer).Value)
		if args[2].Type() != INTEGER_OBJ {
			return NewError("argument to `end` must be INTEGER, got %s", args[2].Type())
		}
		end = int(args[2].(*Integer).Value)
		for i := start; i < end; i++ {
			arr.Elements[i] = value
		}
		return arr
	}
	return nil
}

func array_max(arr *Array) Object {
	if len(arr.Elements) == 0 {
		return NewError("max() cannot be called on an empty array")
	}

	intSort, intArr, strSort, strArr := _checkSortType(arr.Elements)

	if intSort {
		max := intArr[0]
		for _, el := range intArr {
			if el > max {
				max = el
			}
		}
		return &Integer{Value: int64(max)}
	}

	if strSort {
		max := strArr[0]
		for _, el := range strArr {
			if el > max {
				max = el
			}
		}
		return &String{Value: max}
	}

	return NewError("max() only supports primitive types (integers and strings)")
}

func array_sort(arr *Array) Object {
	elements := arr.Elements
	intsort, intArr, strsort, strArr := _checkSortType(elements)

	if !intsort && !strsort {
		return NewError("sort method only works on arrays of integers or strings")
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

	return arr
}
