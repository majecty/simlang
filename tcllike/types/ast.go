package types

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

type LinesNode struct {
	Lines []ASTNode
}

type SymbolNode struct {
	Name string
}

type NumberNode struct {
	Value float64
}

type CallNode struct {
	FuncName string
	Args     []ASTNode
}

func (n *LinesNode) astNode()  {}
func (n *SymbolNode) astNode() {}
func (n *NumberNode) astNode() {}
func (n *CallNode) astNode()   {}

func (n *LinesNode) String() string {
	// TODO: pring line number ai!
	lines := make([]string, len(n.Lines))
	for i, line := range n.Lines {
		lines[i] = line.String()
	}
	return strings.Join(lines, "\n")
}

func (n *SymbolNode) String() string {
	return n.Name
}

func (n *NumberNode) String() string {
	return fmt.Sprintf("%f", n.Value)
}

func (n *CallNode) String() string {
	args := make([]string, len(n.Args))
	for i, arg := range n.Args {
		args[i] = arg.String()
	}
	return fmt.Sprintf("%s(%s)", n.FuncName, strings.Join(args, ", "))
}
