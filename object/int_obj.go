package object

import "strconv"

func intInvokables(method string, i *Integer) Object {
	if method == "to_string" {
		stringVal := strconv.FormatInt(i.Value, 10)
		return &String{Value: stringVal}
	}
	return nil
}
