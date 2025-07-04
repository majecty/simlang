package main

import (
	"fmt"
	"strings"
)

type ASTNode interface {
	astNode()
	String() string
}

type AST struct {
	Root ASTNode
}

type NumberNode struct {
	Value float64
}

type SymbolNode struct {
	Name string
}

type CallNode struct {
	Function ASTNode
	Args     []ASTNode
}

type LetNode struct {
	LetEnv map[string]ASTNode
  Body   ASTNode
}

func (n *NumberNode) astNode() {}
func (n *SymbolNode) astNode() {}
func (n *CallNode) astNode() {}
func (n *LetNode) astNode() {}

func (n *NumberNode) String() string {
	return fmt.Sprintf("Number(%f)", n.Value)
}
func (n *SymbolNode) String() string {
	return fmt.Sprintf("Symbol(%s)", n.Name)
}

func (n *CallNode) String() string {
	args := make([]string, len(n.Args))
	for i, arg := range n.Args {
		args[i] = arg.String()
	}
	return fmt.Sprintf("Call(%s, %s)", n.Function, strings.Join(args, ", "))
}

func (n *CallNode) Push(arg ASTNode) {
	if n.Function == nil {
		n.Function = arg
	} else {
		n.Args = append(n.Args, arg)
	}
}

func (n *LetNode) String() string {
  return fmt.Sprintf("Let(%s, %s)", n.LetEnv, n.Body.String())
}
