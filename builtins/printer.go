package builtins

import (
	"esolang/lang-esolang/object"
	"fmt"
	"strings"
)

func Print(args ...object.Object) object.Object {
	toBePrinted := []string{}
	for _, arg := range args {
		fmt.Print(arg.Inspect())
		toBePrinted = append(toBePrinted, arg.Inspect())
	}
	return &object.String{Value: strings.Join(toBePrinted, " ") + "flag=noshow"}
}

func Println(args ...object.Object) object.Object {
	toBePrinted := []string{}
	for _, arg := range args {
		toBePrinted = append(toBePrinted, arg.Inspect())
		fmt.Println(arg.Inspect())
	}
	return &object.String{Value: strings.Join(toBePrinted, "\n") + "flag=noshow"}
}

// TODO: Rethink how read should work
func Read() object.Object {
	var input string
	fmt.Scanln(&input)
	return &object.String{Value: input}
}
