package main

import (
	"simlang/util"
	"strconv"
)

func parse(tokens []Token) AST {
	var ast AST
	var stack []*CallNode
	initialEmptyCallNode := CallNode{}
	var currentCall *CallNode = &initialEmptyCallNode

	for _, token := range tokens {
		switch token.Type {
		case LPAREN:
			stack = append(stack, currentCall)
			currentCall = &CallNode{}
		case RPAREN:
			util.Invariant(len(stack) > 0, "unbalanced parentheses %v", tokens)
			if stack[len(stack)-1] != &initialEmptyCallNode {
				last := currentCall
				currentCall = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				currentCall.Args = append(currentCall.Args, last)
			}
		case ATOM:
			currentCall.Push(&SymbolNode{Name: token.Value})
		case NUMBER:
			currentCall.Push(&NumberNode{Value: parseFloat64(token.Value)})
		}
	}

	util.Invariant(len(stack) == 1, "unbalanced parentheses %v, len(stack)=%d", tokens, len(stack))
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
