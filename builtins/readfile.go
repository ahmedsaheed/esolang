package builtins

import (
	"esolang/lang-esolang/object"
	"os"
)

/*
ReadFile reads a file and returns the content of the file as a string

	@param args ...object.Object
	@exception wrong number of arguments. got=%d, want=1
	@exception argument to `readFile` must be STRING, got %s
	@exception IOError: Eror reading file %s
*/
func ReadFile(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("Invalid Arguement got=%d, expected=1", len(args))
	}

	if args[0].Type() != object.STRING_OBJ {
		return newError("Supplied path is must be of type String, got %s", args[0].Type())
	}

	file := args[0].(*object.String).Value
	inputFile, err := os.ReadFile(file)

	// check if file exists
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return newError("I/O Error: File %s does not exist", file)
	}

	if err != nil {
		return newError("I/O Error: Error reading file %s", file)
	}

	return &object.String{Value: string(inputFile)}
}
