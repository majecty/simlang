// Package evaluator implements the evaluation logic for the Simlang programming language.
package evaluator

import (
	"fmt"

	"simlang/tcllike/types"
)

func Eval(ast *types.AST) error {
	lines := ast.Root

	if _, err := evalLines(lines); err != nil {
		return fmt.Errorf("failed to eval lines: %w", err)
	}
	return nil
}

func evalLines(lines *types.LinesNode) (any, error) {
	var lastValue any
	lastValue = nil
	for _, line := range lines.Lines {
		if value, err := evalLine(line); err != nil {
			return nil, fmt.Errorf("failed to eval line: %w", err)
		} else {
			lastValue = value
		}
	}
	return lastValue, nil
}

func evalLine(line types.ASTNode) (any, error) {
	switch v := line.(type) {
	case *types.CallNode:
		return evalCall(v)
	case *types.NumberNode:
		return v.Value, nil
	default:
		return nil, fmt.Errorf("not implemented yet for type %T", line)
	}
}

func evalCall(call *types.CallNode) (any, error) {
	f := call.FuncName

	switch f {
	case "print":
		return evalPrint(call.Args)
	case "exp":
		return evalExp(call.Args)
	default:
		return nil, fmt.Errorf("not implemented yet for function %s", f)
	}
}

func evalPrint(args []types.ASTNode) (any, error) {
	values := make([]any, 0)
	for _, arg := range args {
		value, err := evalValue(arg)
		if err != nil {
			return nil, fmt.Errorf("failed to eval print: %w", err)
		}
		values = append(values, value)
	}

	fmt.Println(values...)
	return nil, nil
}

func evalExp(args []types.ASTNode) (any, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("exp takes 3 arguments, got %d", len(args))
	}

	leftArg, err := evalValue(args[0])
	if err != nil {
		return nil, fmt.Errorf("failed to eval exp(left arg): %w", err)
	}
	leftNum, ok := leftArg.(float64)
	if !ok {
		return nil, fmt.Errorf("failed to eval exp(left arg): expected number, got %T", leftArg)
	}

	operatorNode, ok := args[1].(*types.SymbolNode)
	if !ok {
		return nil, fmt.Errorf("failed to eval exp(operator): expected symbol node, got %T", args[1])
	}

	rightArg, err := evalValue(args[2])
	if err != nil {
		return nil, fmt.Errorf("failed to eval exp(right arg): %w", err)
	}
	rightNum, ok := rightArg.(float64)
	if !ok {
		return nil, fmt.Errorf("failed to eval exp(right arg): expected number, got %T", rightArg)
	}

	switch operatorNode.Name {
	case "+":
		return leftNum + rightNum, nil
	case "-":
		return leftNum - rightNum, nil
	case "*":
		return leftNum * rightNum, nil
	case "/":
		return leftNum / rightNum, nil
	default:
		return nil, fmt.Errorf("not implemented yet for operator %s", operatorNode.Name)
	}
}

func evalValue(arg types.ASTNode) (any, error) {
	switch v := arg.(type) {
	case *types.CallNode:
		return evalCall(v)
	case *types.NumberNode:
		return v.Value, nil
	default:
		return nil, fmt.Errorf("not implemented yet for type %T", arg)
	}
}
