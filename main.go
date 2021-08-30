package main

import (
	"fmt"
	"os/user"
)

func main() {
	x, _ := user.Current()
	fmt.Printf("Hello %s This is Monkey programming language \n", x.Username)
}
