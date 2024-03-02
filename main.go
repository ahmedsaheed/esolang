package main

import (
	"fmt"
	"monkey/lang-monkey/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Howdy %s! Welcome to a tiny esolang repl \n", user.Username)
	fmt.Printf("Feel free to type in commands \n")
	repl.Start(os.Stdin, os.Stdout)
}
