package object

func arrayInvokables(method string, arr *Array) Object {
	if method == "count" || method == "length" {
		return &Integer{Value: int64(len(arr.Elements))}
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
