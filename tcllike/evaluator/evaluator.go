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
	case "+":
		return evalAdd(call.Args)
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

func evalAdd(args []types.ASTNode) (any, error) {
	var sum float64 = 0
	for _, arg := range args {
		value, err := evalValue(arg)
		if err != nil {
			return nil, fmt.Errorf("failed to eval print: %w", err)
		}
		sum += value.(float64)
	}
	return sum, nil
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
