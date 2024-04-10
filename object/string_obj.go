package object

import (
	"strconv"
	"unicode/utf8"
)

func stringInvokables(method string, s *String) Object {
	if method == "count" || method == "length" {
		return &Integer{Value: int64(utf8.RuneCountInString(s.Value))}
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
