package object

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

func _noArgsExpected(fnName string, args ...Object) error {
	err := CheckTypings(
		fnName, args,
		ExactArgsLength(0),
	)
	return err
}

func stringInvokables(method string, s *String, args ...Object) Object {
	name := "String." + method
	switch method {
	case "count", "length":
		if err := _noArgsExpected(name, args...); err != nil {
			return newErrorFromTypings(err.Error())
		}
		return &Integer{Value: int64(utf8.RuneCountInString(s.Value))}

	case "reverse":
		if err := _noArgsExpected(name, args...); err != nil {
			return newErrorFromTypings(err.Error())
		}
		runes := []rune(s.Value)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return &String{Value: string(runes)}

	case "equals":
		if err := CheckTypings(
			name, args,
			ExactArgsLength(1),
			WithTypes(STRING_OBJ),
		); err != nil {
			return newErrorFromTypings(err.Error())
		}
		return &Boolean{Value: s.Value == args[0].(*String).Value}

	case "empty":
		if err := _noArgsExpected(name, args...); err != nil {
			return newErrorFromTypings(err.Error())
		}
		return &Boolean{Value: utf8.RuneCountInString(s.Value) == 0}
	case "upper_case":
		if err := _noArgsExpected(name, args...); err != nil {
			return newErrorFromTypings(err.Error())
		}
		return &String{Value: strings.ToUpper(s.Value)}

	case "lower_case":
		if err := _noArgsExpected(name, args...); err != nil {
			return newErrorFromTypings(err.Error())
		}
		return &String{Value: strings.ToLower(s.Value)}

	case "to_int":
		if err := _noArgsExpected(name, args...); err != nil {
			return newErrorFromTypings(err.Error())
		}
		i, err := strconv.ParseInt(s.Value, 0, 64)
		if err != nil {
			i = 0
		}
		return &Integer{Value: int64(i)}
	default:
		return nil
	}
}
