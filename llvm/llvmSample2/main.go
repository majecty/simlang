package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
	"simlang/lexer"
	"simlang/parser"
)

// https://llvm.org/docs/LangRef.html
func main() {
	fmt.Println("Hello, Go Project!")

	if result, err := parser.Parse(lexer.Toknize("(+ 1 2)")); err != nil {
		log.Fatalf("failed to parse %v", err)
	} else {
		fmt.Printf("parse result %v\n", result)
	}

	functionBody := `
	%1 = add i32 0, 1
	%2 = add i32 0, 2
	%3 = add i32 %1, %2

	ret i32 %3
	`
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

@pat = global [14 x i8] c"answer is %d\0A\00"

define i32 @foo() {
	{{.functionBody}}
;	%1 = add i32 0, 1
;	%2 = add i32 0, 2
;	%3 = add i32 %1, %2

;	ret i32 %3
}

define i32 @main() {
	%4 = call i32 @foo()
  %ret = call i32 (ptr, ...) @printf(ptr @pat, i32 %4)
  ret i32 0
}`

	t, err := template.New("llvmTemplate").Parse(llvmIRTemplate)
  if err != nil {
    return "", fmt.Errorf("failed to create template: %v", err)
  }
	var llvmIRSB strings.Builder;
  err = t.Execute(&llvmIRSB, map[string]any{
    "functionBody": functionBody,
  })
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %v", err)
  }
  llvmIR := llvmIRSB.String()
	return llvmIR, nil
}

func writeToFile(filename string, content string) (err error) {
	// Create or open the file for writing
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", filename, err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			errors.Join(err, fmt.Errorf("error closing file %s: %v", filename, closeErr))
		}
	}()

	// Write the LLVM IR string to the file
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write LLVM IR to file %s: %v", filename, err)
	}
	return nil
}

