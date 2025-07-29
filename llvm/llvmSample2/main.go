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
	"simlang/types"
)

// https://llvm.org/docs/LangRef.html
func main() {
	fmt.Println("Hello, Go Project!")

	ast, err := parser.Parse(lexer.Toknize("(+ 1 2)"))
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

type IRGenerationContext struct {
	nextVariableIndex uint32
	instructions      []string
}

func (c *IRGenerationContext) astToLLVMIR(ast *types.AST) (string, error) {
	root := ast.Root
	return c.nodeToLLVMIR(root)
}

func (c *IRGenerationContext) nodeToLLVMIR(node types.ASTNode) (string, error) {
	switch v := node.(type) {
	case *types.CallNode:
		fmt.Printf("root is call node %v\n", v)

		if _, ok := v.Function.(*types.SymbolNode); !ok {
			return "", fmt.Errorf("function is not symbol")
		}
		if v.Function.(*types.SymbolNode).Name != "+" {
			return "", fmt.Errorf("function name is not +")
		}

		if len(v.Args) != 2 {
			return "", fmt.Errorf("args count is not 2")
		}

		arg0, arg0Err := c.nodeToLLVMIRValue(v.Args[0])
		if arg0Err != nil {
			return "", fmt.Errorf("failed to nodeToLLVMIRValue arg0: %w", arg0Err)
		}

		arg1, arg1Err := c.nodeToLLVMIRValue(v.Args[1])
		if arg1Err != nil {
			return "", fmt.Errorf("failed to nodeToLLVMIRValue arg1: %w", arg1Err)
		}

		tempName := c.PutAddInstruction(arg0, arg1)
		c.PutReturnInstruction(tempName)
		return c.MakeFunctionBody(), nil
	default:
		return "", fmt.Errorf("not implemented yet")
	}
}

type IRValue interface {
	toIRValue() string
}

type NumberLiteral struct {
	Value int32
}

func (n *NumberLiteral) toIRValue() string {
	return fmt.Sprintf("%d", n.Value)
}

type RegisterName struct {
	Name string
}

func (r *RegisterName) toIRValue() string {
	return r.Name
}

func (c *IRGenerationContext) nodeToLLVMIRValue(node types.ASTNode) (IRValue, error) {
	switch v := node.(type) {
	case *types.NumberNode:
		log.Printf("nodeToLLVMIRValue NumberNode %v\n", v)
		return &NumberLiteral{Value: int32(v.Value)}, nil
	case *types.SymbolNode:
		return nil, fmt.Errorf("nodeToLLVMIRValue NumberNode not implemented yet")
	case *types.CallNode:
		return nil, fmt.Errorf("nodeToLLVMIRValue CallNode not implemented yet")
	}

	return nil, fmt.Errorf("not implemented yet")
}

/* PutAddInstruction returns new variable name */
func (c *IRGenerationContext) PutAddInstruction(arg0 IRValue, arg1 IRValue) string {
	varName := c.AcquireNextTempVariableName()
	c.instructions = append(c.instructions, fmt.Sprintf("%s = add i32 %s, %s", varName, arg0.toIRValue(), arg1.toIRValue()))
	return varName
}

func (c *IRGenerationContext) AcquireNextTempVariableName() string {
	name := fmt.Sprintf("%%temp.%d", c.nextVariableIndex)
	c.nextVariableIndex += 1
	return name
}

func (c *IRGenerationContext) PutReturnInstruction(varName string) {
	c.instructions = append(c.instructions, fmt.Sprintf("ret i32 %s", varName))
}

func (c *IRGenerationContext) MakeFunctionBody() string {
	stringbuilder := strings.Builder{}
	for _, instruction := range c.instructions {
		stringbuilder.WriteString(fmt.Sprintf("%s\n", instruction))
	}
	return stringbuilder.String()
}
