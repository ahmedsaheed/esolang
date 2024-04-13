package object

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

func stringInvokables(method string, s *String) Object {
	if method == "count" || method == "length" {
		return &Integer{Value: int64(utf8.RuneCountInString(s.Value))}
	}
	if method == "reverse" {
		runes := []rune(s.Value)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return &String{Value: string(runes)}
	}
	if method == "empty" {
		return &Boolean{Value: utf8.RuneCountInString(s.Value) == 0}
	}

	if method == "upper_case" {
		return &String{Value: strings.ToUpper(s.Value)}
	}

	if method == "lower_case" {
		return &String{Value: strings.ToLower(s.Value)}
	}

	if method == "to_int" {
		i, err := strconv.ParseInt(s.Value, 0, 64)
		if err != nil {
			i = 0
		}
		return &Integer{Value: int64(i)}
	}
	return nil
}
