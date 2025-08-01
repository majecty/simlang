package main

import (
	"fmt"
	"log"
	"strings"

	"simlang/types"
)

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
		tempName, err := c.nodeToLLVMIRValue(v)
		if err != nil {
			return "", fmt.Errorf("failed to nodeToLLVMIRValue: %w", err)
		}

		c.PutReturnInstruction(tempName)
		return c.MakeFunctionBody(), nil
	case *types.LetNode:
		tempName, err := c.nodeToLLVMIRValue(v)
		if err != nil {
			return "", fmt.Errorf("failed to nodeToLLVMIRValue: %w", err)
		}

		c.PutReturnInstruction(tempName)
		return c.MakeFunctionBody(), nil
	default:
		return "", fmt.Errorf("not implemented yet for type %T", node)
	}
}

func (c *IRGenerationContext) nodeToLLVMIRValue(node types.ASTNode) (IRValue, error) {
	switch v := node.(type) {
	case *types.NumberNode:
		log.Printf("nodeToLLVMIRValue NumberNode %v\n", v)
		return &NumberLiteral{Value: v.Value}, nil

	case *types.SymbolNode:
		return nil, fmt.Errorf("nodeToLLVMIRValue SymbolNode not implemented yet")

	case *types.CallNode:
		return c.callNodeToLLVMIRValue(v)

	case *types.LetNode:
		return nil, fmt.Errorf("letnode not implemented yet %v", node)
	}

	return nil, fmt.Errorf("not implemented yet %v", node)
}

/* PutFAddInstruction returns new variable name */
func (c *IRGenerationContext) PutFAddInstruction(arg0 IRValue, arg1 IRValue) RegisterName {
	varName := c.AcquireNextTempVariableName()
	c.instructions = append(c.instructions, fmt.Sprintf("%s = fadd double %s, %s", varName, arg0.toIRValue(), arg1.toIRValue()))
	return RegisterName{Name: varName}
}

func (c *IRGenerationContext) AcquireNextTempVariableName() string {
	name := fmt.Sprintf("%%temp.%d", c.nextVariableIndex)
	c.nextVariableIndex += 1
	return name
}

func (c *IRGenerationContext) PutReturnInstruction(irValue IRValue) {
	c.instructions = append(c.instructions, fmt.Sprintf("ret double %s", irValue.toIRValue()))
}

func (c *IRGenerationContext) MakeFunctionBody() string {
	stringbuilder := strings.Builder{}
	for _, instruction := range c.instructions {
		stringbuilder.WriteString(fmt.Sprintf("%s\n", instruction))
	}
	return stringbuilder.String()
}

func (c *IRGenerationContext) callNodeToLLVMIRValue(callNode *types.CallNode) (IRValue, error) {
	if _, ok := callNode.Function.(*types.SymbolNode); !ok {
		return nil, fmt.Errorf("function is not symbol")
	}
	if callNode.Function.(*types.SymbolNode).Name != "+" {
		return nil, fmt.Errorf("function name is not +")
	}

	if len(callNode.Args) == 0 {
		return nil, fmt.Errorf("the number of arguments should be greater than zero")
	}

	sumName, arg0Err := c.nodeToLLVMIRValue(callNode.Args[0])
	if arg0Err != nil {
		return nil, fmt.Errorf("failed to generate ir from arg0 %w", arg0Err)
	}

	for index, arg := range callNode.Args {
		if index == 0 {
			continue
		}

		argI, argErr := c.nodeToLLVMIRValue(arg)
		if argErr != nil {
			return nil, fmt.Errorf("failed to nodeToLLVMIRValue arg: %w", argErr)
		}

		addResult := c.PutFAddInstruction(sumName, argI)
		sumName = &addResult
	}

	return sumName, nil
}
