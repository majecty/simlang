package main

import (
	"fmt"
	"log"
	"strings"

	"simlang/types"
)

type IRRegisterLookup struct {
	prev *IRRegisterLookup
	dict map[string]IRValue
}

type IRGenerationContext struct {
	nextVariableIndex uint32
	instructions      []string
	lookup            *IRRegisterLookup
}

func NewIRGenerationContext() *IRGenerationContext {
	return &IRGenerationContext{
		nextVariableIndex: 0,
		instructions:      []string{},
		lookup: &IRRegisterLookup{
			prev: nil,
			dict: map[string]IRValue{},
		},
	}
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
		irValue := c.lookup.dict[v.Name]
		if irValue != nil {
			return irValue, nil
		}
		return nil, fmt.Errorf("symbol %s not found", v.Name)

	case *types.CallNode:
		return c.callNodeToLLVMIRValue(v)

	case *types.LetNode:
		c.pushLookup()
		defer c.popLookup()
		for argName, arg := range v.LetEnv {
			irValue, err := c.nodeToLLVMIRValue(arg)
			if err != nil {
				return nil, fmt.Errorf("failed to nodeToLLVMIRValue arg: %w", err)
			}
			if err := c.PutLookup(argName, irValue); err != nil {
				return nil, fmt.Errorf("failed to PutLookup: %w", err)
			}
		}

		return c.nodeToLLVMIRValue(v.Body)
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

func (c *IRGenerationContext) PutLookup(name string, irValue IRValue) error {
	if c.lookup == nil {
		return fmt.Errorf("lookup is nil, while putting %s (%v) to lookup", name, irValue)
	}

	if _, ok := c.lookup.dict[name]; ok {
		return fmt.Errorf("lookup already has %s", name)
	}

	c.lookup.dict[name] = irValue
	return nil
}

func (c *IRGenerationContext) pushLookup() {
	c.lookup = &IRRegisterLookup{
		prev: c.lookup,
		dict: map[string]IRValue{},
	}
}

func (c *IRGenerationContext) popLookup() {
	c.lookup = c.lookup.prev
}
