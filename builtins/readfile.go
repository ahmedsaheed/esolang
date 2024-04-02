package builtins

import (
	"esolang/lang-esolang/object"
	"os"
)

/*
readFile reads a file and returns the content of the file as a string

	@param fileLocation string path to the file
	@exception wrong number of arguments.
	@exception wrong type of arguments.
	@exception I/O Error: possible errors reading file
*/
func readFile(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("Invalid Arguement got=%d, expected=1", len(args))
	}

	if args[0].Type() != object.STRING_OBJ {
		return newError("Supplied path must be of type String, got %s", args[0].Type())
	}

	file := args[0].(*object.String).Value
	inputFile, err := os.ReadFile(file)

	// check if file exists
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return newError("I/O Error: No such file %s", file)
	}

	if err != nil {
		return newError("I/O Error: Error reading file %s", file)
	}

	return &object.String{Value: string(inputFile)}
}
