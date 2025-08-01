package main

import (
	"fmt"

	"simlang/tcllike/lexer"
	"simlang/tcllike/parser"
)

func main() {
	fmt.Println("Hello, Go Project!")

	fmt.Println(lexer.Tokenize("print 3"))
	fmt.Println(lexer.Tokenize("print (exp 1 + 2)"))

	fmt.Println(parser.Parse(lexer.Tokenize("print 3")))
	fmt.Println(parser.Parse(lexer.Tokenize("print (exp 1 + 2)")))

	fmt.Println(parser.Parse(lexer.Tokenize(`print (exp 1 + 2)
		print 3`)))
}
