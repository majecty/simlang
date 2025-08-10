package main

import (
	"fmt"
	"strings"

	"simlang/tcllike/types"
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
	case *types.LinesNode:
		var lastIRValue IRValue = nil
		for _, line := range v.Lines {
			lineResult, err := c.nodeToLLVMIRValue(line)
			if err != nil {
				return "", fmt.Errorf("failed to nodeToLLVMIR: %w", err)
			}
			lastIRValue = lineResult
		}

		c.PutReturnInstruction(lastIRValue)
		return c.MakeFunctionBody(), nil
	default:
		return "", fmt.Errorf("not implemented yet for type %T", node)
	}
}

func (c *IRGenerationContext) nodeToLLVMIRValue(node types.ASTNode) (IRValue, error) {
	switch v := node.(type) {
	case *types.NumberNode:
		return &NumberLiteral{Value: v.Value}, nil
	}
	return nil, fmt.Errorf("not implemented yet for type %T", node)
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
