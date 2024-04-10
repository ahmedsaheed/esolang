package main

import _ "embed"

//go:embed stdlib/array.eso
var ArrayUtils string

//go:embed stdlib/bool.eso
var BoolUtils string

//go:embed stdlib/string.eso
var StringUtils string

//go:embed stdlib/set.eso
var SetUtils string

func getStdLib() string {
	return ArrayUtils + "\n" + BoolUtils + "\n" + StringUtils + "\n" + SetUtils

}
