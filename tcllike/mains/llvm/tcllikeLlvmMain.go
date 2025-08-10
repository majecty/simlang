package main

import (
	"log"

	"simlang/tcllike/lexer"
	"simlang/tcllike/parser"
)

func main() {
	println("Hello, World!")

	ast, err := parser.Parse(lexer.Tokenize("(3 + 4)"))
	if err != nil {
		panic(err)
	}

	iRGenerationContext := IRGenerationContext{}
	body, irErr := iRGenerationContext.astToLLVMIR(ast)
	if irErr != nil {
		log.Fatalf("failed to astToLLVMIR %v", irErr)
	}
	println(body)
}
