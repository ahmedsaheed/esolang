package builtins

import (
	"esolang/lang-esolang/object"
	"fmt"
	"net/http"
)

func _httpServeAndListen(args ...object.Object) object.Object {
	// param 1: port
	// param 2: route
	// param 3: response

	port := args[0].(*object.Integer).Value
	route := args[1].(*object.String).Value
	response := args[2].(*object.String).Value

	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, response)
	}

	http.HandleFunc(route, handler)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	return NULL
}
