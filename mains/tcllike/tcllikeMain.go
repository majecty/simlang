package main

import (
	"fmt"

	"simlang/tcllike/lexer"
)

func main() {
	fmt.Println("Hello, Go Project!")

	fmt.Println(lexer.Tokenize("print 3"))
}
