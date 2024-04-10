package object

func intInvokables(method string, i *Integer) Object {
	if method == "to_string" {
		return &String{Value: string(rune(i.Value))}
	}
	return nil
}
