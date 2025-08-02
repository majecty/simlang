package main

import (
	"fmt"

	"simlang/tcllike/evaluator"
	"simlang/tcllike/lexer"
	"simlang/tcllike/parser"
)

func main() {
	fmt.Println("Hello, Go Project!")
	lpe(`print 3`)
	lpe(`print (1 + (2 + 3))`)
	lpe(`print [+ 1 2]`)
	lpe(`print (1 + 2)
    print 3`)
}

func lpe(code string) {
	tokens := lexer.Tokenize(code)
	ast, astErr := parser.Parse(tokens)
	if astErr != nil {
		panic(astErr)
	}
	fmt.Printf("%v => \n", code)
	evaluator.Eval(ast)
}
