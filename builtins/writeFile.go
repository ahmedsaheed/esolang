package builtins

import (
	"esolang/lang-esolang/object"
	"os"
)

/*
writeFile writes data to a named file, creating it if necessary

	@param args ...object.Object
	@expectedParam args[0] string path to the file
	@expectedParam args[1] string content to write to the file
	@exception wrong number of arguments.
	@exception wrong type of arguments.
	@exception I/O Error: possible error writing to file
*/
func writeFile(arg ...object.Object) object.Object {
	if len(arg) < 2 || len(arg) > 3 {
		return newError("Invalid Arguement got=%d, expected=3", len(arg))
	}

	if arg[0].Type() != object.STRING_OBJ {
		return newError("Supplied path must be of type String, got %s", arg[0].Type())
	}

	if arg[1].Type() != object.STRING_OBJ {
		return newError("Supplied content must be of type String, got %s", arg[1].Type())
	}

	// check for flag - append or prepend
	if len(arg) == 3 {
		if arg[2].Type() != object.STRING_OBJ {
			return newError("Supplied flag must be of type String, got %s", arg[2].Type())
		}
		flag := arg[2].(*object.String).Value

		// warn invalid flag
		if flag != "a+" && flag != "+a" {
			return newError("Invalid flag for writeFile %s, expected a+ or +a", flag)
		}

		// append to file
		if flag == "a+" {
			return appendToFile(arg[0].(*object.String).Value, arg[1].(*object.String).Value)
		}

		// prepend to file
		if flag == "+a" {
			return prependToFile(arg[0].(*object.String).Value, arg[1].(*object.String).Value)
		}
	}

	file := arg[0].(*object.String).Value
	content := arg[1].(*object.String).Value

	// create the file if it does not exist with 0666 perm
	err := os.WriteFile(file, []byte(content), 0666)

	if err != nil {
		return newError("I/O Error: Error writing file %s", file)
	}

	return NULL
}

func appendToFile(filePath string, content string) object.Object {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return newError("I/O Error: Error opening to file %s", filePath)
	}
	defer f.Close()
	data := []byte(content)
	if _, err = f.Write(data); err != nil {
		return newError("I/O Error: Error writing to file %s", filePath)
	}
	return NULL
}

func prependToFile(filePath string, content string) object.Object {
	existingData, err := os.ReadFile(filePath)
	if err != nil {
		return newError("I/O Error: Error reading file %s", filePath)
	}
	newData := []byte(content + "\n")
	newData = append(newData, existingData...)

	err = os.WriteFile(filePath, newData, 0666)
	return NULL
}
