package main

import (
	"fmt"
	"monkey/repl"
	"os"
	"os/user"
)

func main() {
	x, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s This is Monkey programming language \n", x.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}
