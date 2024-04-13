package object

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

func array_fill(method string, arr *Array, args ...Object) Object {
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
			return newError("argument to `start` must be INTEGER, got %s", args[1].Type())
		}
		start = int(args[1].(*Integer).Value)
		for i := start; i < len(arr.Elements); i++ {
			arr.Elements[i] = value
		}
		return arr
	}

	if len(args) == 3 {
		if args[1].Type() != INTEGER_OBJ {
			return newError("argument to `start` must be INTEGER, got %s", args[1].Type())
		}
		start = int(args[1].(*Integer).Value)
		if args[2].Type() != INTEGER_OBJ {
			return newError("argument to `end` must be INTEGER, got %s", args[2].Type())
		}
		end = int(args[2].(*Integer).Value)
		for i := start; i < end; i++ {
			arr.Elements[i] = value
		}
		return arr
	}
	return nil
}
