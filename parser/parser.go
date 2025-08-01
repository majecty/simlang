// Package parser implements the parsing logic for the Simlang programming language.
// It transforms tokens from the lexer into an abstract syntax tree (AST) representation.
// The parser handles language constructs like function calls, let expressions, and lambdas.
package parser

import (
	"errors"
	"fmt"
	"strconv"

	"simlang/types"
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

func (p *ParsingContext) consume() types.Token {
	if p.currentTokenIndex >= len(p.tokens) {
		panic("unexpected end of input in consume")
	}

	token := p.currentToken()
	p.currentTokenIndex++
	return token
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
	node, err := parseSingle(&parsingContext)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %v: %w", tokens, err)
	}
	if parsingContext.hasNextToken() {
		return nil, fmt.Errorf("node is parsed but tokens remains, node: %v", node)
	}

	return &types.AST{Root: node}, nil
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

func parseSingle(parsingContext *ParsingContext) (types.ASTNode, error) {
	switch parsingContext.currentToken().Type {
	case types.LPAREN:
		return parseFromLParen(parsingContext)
	case types.ATOM:
		return &types.SymbolNode{Name: parsingContext.consume().Value}, nil
	case types.NUMBER:
		return &types.NumberNode{Value: parseFloat64(parsingContext.consume().Value)}, nil
	default:
		return nil, fmt.Errorf("in parseSingle, unexpected main.TokenType %v", parsingContext.currentToken().Type.String())
	}
}

func parseFromLParen(parsingContext *ParsingContext) (types.ASTNode, error) {
	token := parsingContext.consume()
	util.Invariant(token.Type == types.LPAREN, "invalid atom, there should be lparentheses %v", parsingContext.tokens)

	token = parsingContext.consume()
	switch token.Type {
	case types.ATOM:
		parsingContext.back()
		parsingContext.back()
		funcCallNode, err := parseFunctionCall(parsingContext)
		if err != nil {
			return nil, fmt.Errorf("failed to parse function call: %w", err)
		}
		return funcCallNode, nil
	case types.LET:
		parsingContext.back()
		parsingContext.back()
		letValueNode, err := parseLetValue(parsingContext)
		if err != nil {
			return nil, fmt.Errorf("failed to parse let value: %w", err)
		}
		return letValueNode, nil
	case types.LAMBDA:
		parsingContext.back()
		parsingContext.back()
		lambdaNode, err := parseLambda(parsingContext)
		if err != nil {
			return nil, fmt.Errorf("failed to parse lambda: %w", err)
		}
		return lambdaNode, nil
	case types.NUMBER:
		return nil, fmt.Errorf("there should be function call but found number %v", token)
	case types.LPAREN:
		return nil, fmt.Errorf("there should be function call but found lparen")
	case types.RPAREN:
		return nil, fmt.Errorf("invalid rparen, we are not supporting empty list")
	default:
		return nil, fmt.Errorf("unexpected main.TokenType: %#v", token.Type.String())
	}
}

func discardRParen(parsingContext *ParsingContext) error {
	token := parsingContext.consume()
	if token.Type != types.RPAREN {
		return fmt.Errorf("expected rparent but get:  %v", token)
	}
	return nil
}

func parseFunctionCall(parsingContext *ParsingContext) (*types.CallNode, error) {
	if err := discardLParen(parsingContext); err != nil {
		return nil, fmt.Errorf("failed to parse function call: %w", err)
	}

	funcSymbolNode, err := parseSymbol(parsingContext)
	if err != nil {
		return nil, fmt.Errorf("failed to parse function call node: %w", err)
	}

	args := make([]types.ASTNode, 0)
	for parsingContext.currentToken().Type != types.RPAREN {
		argNode, err := parseSingle(parsingContext)
		if err != nil {
			return nil, fmt.Errorf("failed to parse function call arg: %w", err)
		}
		args = append(args, argNode)
	}

	if err := discardRParen(parsingContext); err != nil {
		return nil, fmt.Errorf("failed to parse function call, handling last rparen: %w", err)
	}

	return &types.CallNode{Function: funcSymbolNode, Args: args}, nil
}

func discardLParen(parsingContext *ParsingContext) error {
	token := parsingContext.consume()
	if token.Type != types.LPAREN {
		return fmt.Errorf("expected lparen but get:  %v", token)
	}
	return nil
}

func parseSymbol(parsingContext *ParsingContext) (*types.SymbolNode, error) {
	token := parsingContext.consume()
	if token.Type != types.ATOM {
		return nil, fmt.Errorf("expected atom but get:  %v", token)
	}
	return &types.SymbolNode{Name: token.Value}, nil
}

// (let (x 10) in x)
func parseLetValue(parsingContext *ParsingContext) (types.ASTNode, error) {
	if err := discardLParen(parsingContext); err != nil {
		return nil, fmt.Errorf("failed to parse let value: %w", err)
	}
	if err := discardLet(parsingContext); err != nil {
		return nil, fmt.Errorf("failed to parse let value: %w", err)
	}
	if err := discardLParen(parsingContext); err != nil {
		return nil, fmt.Errorf("failed to parse let value: %w", err)
	}

	envVariableName, err := parseSymbol(parsingContext)
	if err != nil {
		return nil, fmt.Errorf("failed to parse let value, while parsing env variable name: %w", err)
	}

	envVariableValue, err := parseSingle(parsingContext)
	if err != nil {
		return nil, fmt.Errorf("failed to parse let value, while parsing env variable value: %w", err)
	}
	if err := discardRParen(parsingContext); err != nil {
		return nil, fmt.Errorf("failed to parse let value: %w", err)
	}

	if err := discardIN(parsingContext); err != nil {
		return nil, fmt.Errorf("failed to parse let value, try to discard 'in': %w", err)
	} // discardIN(parsingContext)

	body, err := parseSingle(parsingContext)
	if err != nil {
		return nil, fmt.Errorf("failed to parse let value, while parsing body: %w", err)
	}
	if err := discardRParen(parsingContext); err != nil {
		return nil, fmt.Errorf("failed to parse let value, try consume last rparen: %w", err)
	}

	return &types.LetNode{
		LetEnv: map[string]types.ASTNode{envVariableName.Name: envVariableValue},
		Body:   body,
	}, nil
}

func discardLet(parsingContext *ParsingContext) error {
	token := parsingContext.consume()
	if token.Type != types.LET {
		return fmt.Errorf("expected let but get:  %v", token)
	}

	return nil
}

func discardIN(parsingContext *ParsingContext) error {
	token := parsingContext.consume()
	if token.Type != types.IN {
		return fmt.Errorf("expected in but get:  %v", token)
	}

	return nil
}

func parseLambda(parsingContext *ParsingContext) (types.ASTNode, error) {
	if err := discardLParen(parsingContext); err != nil {
		return nil, fmt.Errorf("failed to parse lambda: %w", err)
	}
	if err := discardLAMBDA(parsingContext); err != nil {
		return nil, fmt.Errorf("failed to parse lambda: %w", err)
	}
	if err := discardLParen(parsingContext); err != nil {
		return nil, fmt.Errorf("failed to parse lambda: %w", err)
	}

	args := make([]*types.SymbolNode, 0)
	for parsingContext.currentToken().Type != types.RPAREN {
		envVariableName, err := parseSymbol(parsingContext)
		if err != nil {
			return nil, fmt.Errorf("failed to parse lambda, while parsing env variable name: %w", err)
		}
		args = append(args, envVariableName)
	}
	if err := discardRParen(parsingContext); err != nil {
		return nil, fmt.Errorf("failed to parse lambda: %w", err)
	}

	body, err := parseSingle(parsingContext)
	if err != nil {
		return nil, fmt.Errorf("failed to parse lambda, while parsing body: %w", err)
	}
	if err := discardRParen(parsingContext); err != nil {
		return nil, fmt.Errorf("failed to parse lambda, try consume last rparen: %w", err)
	}

	return &types.LambdaNode{
		Args: args,
		Body: body,
	}, nil
}

func discardLAMBDA(parsingContext *ParsingContext) error {
	token := parsingContext.consume()
	if token.Type != types.LAMBDA {
		return fmt.Errorf("expected lambda but get:  %v", token)
	}

	return nil
}
