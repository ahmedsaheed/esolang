package builtins

import (
	_ "embed"
	"fmt"
)

//go:embed stdlib/array.eso
var ArrayUtils string

//go:embed stdlib/bool.eso
var BoolUtils string

//go:embed stdlib/string.eso
var StringUtils string

//go:embed stdlib/set.eso
var SetUtils string

//go:embed stdlib/math.eso
var Math string

func getAllStdLib() string {
	return ArrayUtils + "\n" + BoolUtils + "\n" + StringUtils + "\n" + SetUtils
}

func GetStdLib(lib string) (string, error) {
	switch lib {
	case "array":
		return ArrayUtils, nil
	case "bool":
		return BoolUtils, nil
	case "string":
		return StringUtils, nil
	case "set":
		return SetUtils, nil
	case "math":
		return Math, nil
	default:
		return "", fmt.Errorf("stdlib: %s not found", lib)
	}
}
