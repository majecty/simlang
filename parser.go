package main

import (
	"strconv"
)

func parse(tokens []Token) AST {
	var ast AST
	var stack []CallNode
	initialEmptyCallNode := CallNode{}
	var currentCall *CallNode = &initialEmptyCallNode

	for _, token := range tokens {
		switch token.Type {
		case LPAREN:
			if currentCall != &initialEmptyCallNode {
				stack = append(stack, *currentCall)
			}
			currentCall = &CallNode{}
		case RPAREN:
			if len(stack) > 0 {
				last := currentCall
				currentCall = &stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				currentCall.Args = append(currentCall.Args, last)
			}
		case ATOM:
			currentCall.Push(&SymbolNode{Name: token.Value})
		case NUMBER:
			currentCall.Push(&NumberNode{Value: parseFloat64(token.Value)})
		}
	}

	ast.Root = currentCall
	return ast
}

func parseFloat64(value string) float64 {
	if value == "" {
		return 0
	}
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		panic(err)
	}
	return f
}
