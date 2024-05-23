package object

func setInvokables(method string, set *Set, args ...Object) Object {
	switch method {
	case "insert":
		return setInsert(set, args...)
	case "delete":
		return setRemove(set, args...)
	case "contains":
		return setContains(set, args...)
	case "size":
		return setSize(set)
	case "clear":
		return setClear(set)
	case "isEmpty":
		return setIsEmpty(set)
	}
	return newError("method not found: %s", method)
}

func setInsert(set *Set, args ...Object) Object {
	name := "Set.add"
	if err := CheckTypings(
		name, args,
		ExactArgsLength(1),
	); err != nil {
		return newErrorFromTypings(err.Error())
	}
	inSet, _ := isObjectInSet(args[0], set)
	if !inSet {
		set.Elements = append(set.Elements, args[0])
	}
	return set
}

func setRemove(set *Set, args ...Object) Object {
	name := "Set.remove"
	if err := CheckTypings(
		name, args,
		ExactArgsLength(1),
	); err != nil {
		return newErrorFromTypings(err.Error())
	}
	inSet, idx := isObjectInSet(args[0], set)
	if inSet {
		set.Elements = append(set.Elements[:idx], set.Elements[idx+1:]...)
	}
	return set
}

func setContains(set *Set, args ...Object) Object {
	name := "Set.contains"
	if err := CheckTypings(
		name, args,
		ExactArgsLength(1),
	); err != nil {
		return newErrorFromTypings(err.Error())
	}
	inSet, _ := isObjectInSet(args[0], set)
	return &Boolean{Value: inSet}
}

func setSize(set *Set, args ...Object) Object {
	name := "Set.size"
	if err := CheckTypings(
		name, args,
		ExactArgsLength(0),
	); err != nil {
		return newErrorFromTypings(err.Error())
	}
	return &Integer{Value: int64(len(set.Elements))}
}

func setIsEmpty(set *Set, args ...Object) Object {
	name := "Set.isEmpty"
	if err := CheckTypings(
		name, args,
		ExactArgsLength(0),
	); err != nil {
		return newErrorFromTypings(err.Error())
	}
	return &Boolean{Value: len(set.Elements) == 0}
}

func setClear(set *Set) Object {
	set.Elements = []Object{}
	return set
}

func isObjectInSet(obj Object, set *Set) (bool, int) {
	for idx, el := range set.Elements {
		if obj.Inspect() == el.Inspect() {
			return true, idx
		}
	}
	return false, -1
}
