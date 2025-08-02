// Package parser implements parsing functionality for the Tcl-like language variant.
// It transforms tokens produced by the lexer into an abstract syntax tree (AST).
package parser

import (
	"errors"
	"fmt"
	"strconv"

	"simlang/tcllike/types"
	"simlang/util"
)

type ParsingContext struct {
	tokens            []types.Token
	currentTokenIndex int
}

func (p *ParsingContext) currentToken() types.Token {
	return p.tokens[p.currentTokenIndex]
}

func (p *ParsingContext) hasNextToken() bool {
	return p.currentTokenIndex < len(p.tokens)
}

func (p *ParsingContext) consume() (types.Token, error) {
	if p.currentTokenIndex >= len(p.tokens) {
		return types.Token{}, errors.New("unexpected end of input in consume")
	}
	token := p.currentToken()
	p.currentTokenIndex++
	return token, nil
}

func (p *ParsingContext) back() {
	if p.currentTokenIndex == 0 {
		panic("unexpected end of input in back")
	}
	p.currentTokenIndex--
}

func Parse(tokens []types.Token) (*types.AST, error) {
	if len(tokens) == 0 {
		return nil, errors.New("empty input")
	}

	parsingContext := ParsingContext{tokens: tokens, currentTokenIndex: 0}
	node, err := parseLines(&parsingContext)
	if err != nil {
		return nil, fmt.Errorf("failed to parse lines: %w\n input was %v", err, tokens)
	}
	if parsingContext.hasNextToken() {
		return nil, fmt.Errorf("node is parsed but tokens remains, node: %v", node)
	}

	return &types.AST{Root: node}, nil
}

func parseLines(parsingContext *ParsingContext) (*types.LinesNode, error) {
	lines := make([]types.ASTNode, 0)

	for parsingContext.hasNextToken() {
		node, err := parseCall(parsingContext)
		if err != nil {
			return nil, fmt.Errorf("failed to parse line: %w", err)
		}
		if err := consumeLineEnd(parsingContext); err != nil {
			return nil, fmt.Errorf("failed to parse line: %w", err)
		} // consumeLineEnd(parsingContext)
		lines = append(lines, node)
	}

	return &types.LinesNode{Lines: lines}, nil
}

func parseCall(parsingContext *ParsingContext) (types.ASTNode, error) {
	funcToken, err := parsingContext.consume()
	if err != nil {
		return nil, fmt.Errorf("failed to parse call(first function name): %w", err)
	}
	util.Invariant(funcToken.Type == types.Atom, "invalid atom, there should be function call atom %v", parsingContext.tokens)
	args := make([]types.ASTNode, 0)

	argIndex := 0
loop:
	for parsingContext.hasNextToken() {
		nextToken, err := parsingContext.consume()
		if err != nil {
			return nil, fmt.Errorf("failed to parse call(arg %d): %w", argIndex, err)
		}
		argIndex++
		switch nextToken.Type {
		case types.Atom:
			args = append(args, &types.SymbolNode{Name: nextToken.Value})
		case types.Number:
			args = append(args, &types.NumberNode{Value: parseFloat64(nextToken.Value)})
		case types.LParen:
			parsingContext.back()
			argNode, err := parseParen(parsingContext)
			if err != nil {
				return nil, fmt.Errorf("failed to parse call(arg %d): %w", argIndex, err)
			}
			args = append(args, argNode)
		default:
			parsingContext.back()
			break loop
		}
	}

	return &types.CallNode{FuncName: funcToken.Value, Args: args}, nil
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

func parseParen(parsingContext *ParsingContext) (types.ASTNode, error) {
	token, err := parsingContext.consume()
	if err != nil {
		return nil, fmt.Errorf("failed to parse lparen: %w", err)
	}
	if token.Type != types.LParen {
		return nil, fmt.Errorf("expected lparen, got %v", token)
	}

	elements := make([]types.ASTNode, 0)
loop:
	for parsingContext.hasNextToken() {
		token, err := parsingContext.consume()
		if err != nil {
			return nil, fmt.Errorf("failed to parse lparen: %w", err)
		}

		switch token.Type {
		case types.RParen:
			break loop
		case types.LParen:
			parsingContext.back()
			inner, err := parseParen(parsingContext)
			if err != nil {
				return nil, fmt.Errorf("failed to parse lparen: %w", err)
			}
			elements = append(elements, inner)
		case types.Atom:
			elements = append(elements, &types.SymbolNode{Name: token.Value})
		case types.Number:
			elements = append(elements, &types.NumberNode{Value: parseFloat64(token.Value)})
		}
	}

	if len(elements) == 0 {
		return nil, errors.New("empty elements")
	}

	if len(elements) != 3 {
		return nil, fmt.Errorf("expected 3 elements in paren, got %d", len(elements))
	}

	leftArg := elements[0]
	operator := elements[1]
	rightArg := elements[2]

	operatorSymbol, ok := operator.(*types.SymbolNode)
	if !ok {
		return nil, fmt.Errorf("expected operator to be symbol node, got %T", operator)
	}

	return &types.CallNode{FuncName: operatorSymbol.Name, Args: []types.ASTNode{leftArg, rightArg}}, nil
}

func consumeRParen(parsingContext *ParsingContext) error {
	token, err := parsingContext.consume()
	if err != nil {
		return fmt.Errorf("failed to parse rparen: %w", err)
	}
	if token.Type != types.RParen {
		return fmt.Errorf("expected rparen, got %v", token)
	}
	return nil
}

func consumeLineEnd(parsingContext *ParsingContext) error {
	// eof 인 경우 line end 없어도 무시
	if !parsingContext.hasNextToken() {
		return nil
	}
	token, err := parsingContext.consume()
	if err != nil {
		return fmt.Errorf("failed to parse line end: %w", err)
	}

	if token.Type != types.LineEnd {
		return fmt.Errorf("expected line end, got %v", token)
	}
	return nil
}
