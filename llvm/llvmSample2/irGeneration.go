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

		tempName := c.PutFAddInstruction(arg0, arg1)
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
	Value float64
}

func (n *NumberLiteral) toIRValue() string {
	return fmt.Sprintf("%f", n.Value)
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

		if len(v.Args) != 2 {
			return nil, fmt.Errorf("args count is not 2")
		}

		arg0, arg0Err := c.nodeToLLVMIRValue(v.Args[0])
		if arg0Err != nil {
			return nil, fmt.Errorf("failed to nodeToLLVMIRValue arg0: %w", arg0Err)
		}

		arg1, arg1Err := c.nodeToLLVMIRValue(v.Args[1])
		if arg1Err != nil {
			return nil, fmt.Errorf("failed to nodeToLLVMIRValue arg1: %w", arg1Err)
		}

		tempName := c.PutFAddInstruction(arg0, arg1)
		return &RegisterName{Name: tempName}, nil
	}

	return nil, fmt.Errorf("not implemented yet")
}

/* PutFAddInstruction returns new variable name */
func (c *IRGenerationContext) PutFAddInstruction(arg0 IRValue, arg1 IRValue) string {
	varName := c.AcquireNextTempVariableName()
	c.instructions = append(c.instructions, fmt.Sprintf("%s = fadd double %s, %s", varName, arg0.toIRValue(), arg1.toIRValue()))
	return varName
}

func (c *IRGenerationContext) AcquireNextTempVariableName() string {
	name := fmt.Sprintf("%%temp.%d", c.nextVariableIndex)
	c.nextVariableIndex += 1
	return name
}

func (c *IRGenerationContext) PutReturnInstruction(varName string) {
	c.instructions = append(c.instructions, fmt.Sprintf("ret double %s", varName))
}

func (c *IRGenerationContext) MakeFunctionBody() string {
	stringbuilder := strings.Builder{}
	for _, instruction := range c.instructions {
		stringbuilder.WriteString(fmt.Sprintf("%s\n", instruction))
	}
	return stringbuilder.String()
}
