package main

import (
	"fmt"
	"html/template"
	"log"
	"strings"

	"simlang/lexer"
	"simlang/parser"
)

// https://llvm.org/docs/LangRef.html
func main() {
	fmt.Println("Hello, Go Project!")

	ast, err := parser.Parse(lexer.Toknize("(let (x 10) in x)"))

	if err != nil {
		log.Fatalf("failed to parse %v", err)
	} else {
		fmt.Printf("parse result %v\n", ast)
	}
	var body string
	iRGenerationContext := IRGenerationContext{}
	body, err = iRGenerationContext.astToLLVMIR(ast)
	if err != nil {
		log.Fatalf("failed to astToLLVMIR %v", err)
	}
	log.Printf("generated body is %v", body)

	// functionBody := `
	// %1 = add i32 0, 1
	// %2 = add i32 0, 2
	// %3 = add i32 %1, %2
	//
	// ret i32 %3

	functionBody := body
	llvmIR, err := fillTemplate(functionBody)
	if err != nil {
		log.Fatalf("failed to fill template %v", err)
	}

	filename := "output.ll"

	if err := writeToFile(filename, llvmIR); err != nil {
		log.Fatalf("Failed to write LLVM IR to file %s: %v", filename, err)
	}
	log.Printf("Successfully wrote LLVM IR to %s", filename)
}

func fillTemplate(functionBody string) (string, error) {
	llvmIRTemplate := `; ModuleID = 'simple_module'
source_filename = "simple_program.ll"
declare i32 @printf(ptr, ...)

@pat = global [14 x i8] c"answer is %f\0A\00"

define double @foo() {
	{{.functionBody}}
;	%1 = add i32 0, 1
;	%2 = add i32 0, 2
;	%3 = add i32 %1, %2

;	ret i32 %3
}

define i32 @main() {
	%4 = call double @foo()
  %ret = call i32 (ptr, ...) @printf(ptr @pat, double %4)
  ret i32 0
}`

	t, err := template.New("llvmTemplate").Parse(llvmIRTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to create template: %v", err)
	}
	var llvmIRSB strings.Builder
	err = t.Execute(&llvmIRSB, map[string]any{
		"functionBody": functionBody,
	})
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %v", err)
	}
	llvmIR := llvmIRSB.String()
	return llvmIR, nil
}
