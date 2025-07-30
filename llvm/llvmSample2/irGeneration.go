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
	default:
		return "", fmt.Errorf("not implemented yet")
	}
}

func (c *IRGenerationContext) nodeToLLVMIRValue(node types.ASTNode) (IRValue, error) {
	switch v := node.(type) {
	case *types.NumberNode:
		log.Printf("nodeToLLVMIRValue NumberNode %v\n", v)
		return &NumberLiteral{Value: v.Value}, nil
	case *types.SymbolNode:
		return nil, fmt.Errorf("nodeToLLVMIRValue NumberNode not implemented yet")
	case *types.CallNode:
		if _, ok := v.Function.(*types.SymbolNode); !ok {
			return nil, fmt.Errorf("function is not symbol")
		}
		if v.Function.(*types.SymbolNode).Name != "+" {
			return nil, fmt.Errorf("function name is not +")
		}

		// TODO: 첫번째 arg는 변수 없이 사용도 가능함. 최적화 가능
		sumName := c.PutFAddInstruction(
			&NumberLiteral{Value: 0},
			&NumberLiteral{Value: 0})

		for _, arg := range v.Args {
			argI, argErr := c.nodeToLLVMIRValue(arg)
			if argErr != nil {
				return nil, fmt.Errorf("failed to nodeToLLVMIRValue arg: %w", argErr)
			}

			sumName = c.PutFAddInstruction(&sumName, argI)
		}

		return &sumName, nil
	}

	return nil, fmt.Errorf("not implemented yet")
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
