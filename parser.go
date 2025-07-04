package main

import (
	"errors"
	"fmt"
	"simlang/util"
	"strconv"
)

type ParsingContext struct {
	tokens []Token
	currentTokenIndex int
}

func (p *ParsingContext) currentToken() Token {
  return p.tokens[p.currentTokenIndex]
}

func (p *ParsingContext) hasNextToken() bool {
  return p.currentTokenIndex < len(p.tokens)
}

func (p *ParsingContext) consume() Token {
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

func parse(tokens []Token) (*AST, error) {
	if (len(tokens) == 0) {
		return nil, errors.New("empty input")
	}

	parsingContext := ParsingContext{tokens: tokens, currentTokenIndex: 0}
  node, err := parseSingle(&parsingContext)
	if err != nil {
    return nil, err
	}
	if parsingContext.hasNextToken() {
		return nil, fmt.Errorf("node is parsed but tokens remains, node: %v", node)
	}

	return &AST{Root: node}, nil
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

func parseSingle(parsingContext *ParsingContext) (ASTNode, error) {
	switch parsingContext.currentToken().Type {
	case LPAREN: return parseFromLParen(parsingContext)
	case ATOM:
    return &SymbolNode{Name: parsingContext.consume().Value}, nil
	case LET:
	  return nil, fmt.Errorf("invalid let, there should be lparentheses or value %v")
	case NUMBER:
    return &NumberNode{Value: parseFloat64(parsingContext.consume().Value)}, nil
	case RPAREN:
		return nil, fmt.Errorf("invalid atom, there should be lparentheses or value %v", parsingContext.tokens)
	default:
		panic("unexpected main.TokenType")
	}
}

func parseFromLParen(parsingContext *ParsingContext) (ASTNode, error) {
	token := parsingContext.consume()
	util.Invariant(token.Type == LPAREN, "invalid atom, there should be lparentheses %v", parsingContext.tokens)

	token = parsingContext.consume()
	switch token.Type {
	case ATOM:
		parsingContext.back()
		parsingContext.back()
		funcCallNode, err := parseFunctionCall(parsingContext)
		if err != nil {
			return nil, fmt.Errorf("failed to parse function call: %w", err)
		}
		return funcCallNode, nil
	case LET:
		parsingContext.back()
		parsingContext.back()
		letValueNode, err := parseLetValue(parsingContext)
		if err != nil {
			return nil, fmt.Errorf("failed to parse let value: %w", err)
		}
		return letValueNode, nil
	case NUMBER:
		return nil, fmt.Errorf("there should be function call but found number %v", token)
	case LPAREN:
		return nil, fmt.Errorf("there should be function call but found lparen")
	case RPAREN:
	  return nil, fmt.Errorf("invalid rparen, we are not supporting empty list")
	default:
		panic(fmt.Sprintf("unexpected main.TokenType: %#v", token.Type))
	}
}

func discardRParen(parsingContext *ParsingContext) error {
  token := parsingContext.consume()
	if token.Type != RPAREN {
    return fmt.Errorf("expected rparent but get:  %v", token)
  }
  return nil
}

func parseFunctionCall(parsingContext *ParsingContext) (*CallNode, error) {
	if err := discardLParen(parsingContext); err != nil {
		return nil, fmt.Errorf("failed to parse function call: %w", err)
	}

	funcSymbolNode, err := parseSymbol(parsingContext)
	if err != nil {
    return nil, fmt.Errorf("failed to parse function call node: %w", err)
  }

	args := make([]ASTNode, 0)
	for parsingContext.currentToken().Type != RPAREN {
		argNode, err := parseSingle(parsingContext)
    if err != nil {
      return nil, fmt.Errorf("failed to parse function call arg: %w", err)
    }
    args = append(args, argNode)
	}

	if err := discardRParen(parsingContext); err != nil {
    return nil, fmt.Errorf("failed to parse function call, handling last rparen: %w", err)
  }

	return &CallNode{Function: funcSymbolNode, Args: args}, nil
}

func discardLParen(parsingContext *ParsingContext) error {
  token := parsingContext.consume()
  if token.Type != LPAREN {
    return fmt.Errorf("expected lparen but get:  %v", token)
  }
  return nil
}

func parseSymbol(parsingContext *ParsingContext) (*SymbolNode, error) {
  token := parsingContext.consume()
  if token.Type != ATOM {
    return nil, fmt.Errorf("expected atom but get:  %v", token)
  }
  return &SymbolNode{Name: token.Value}, nil
}

// (let (x 10) in x)
func parseLetValue(parsingContext *ParsingContext) (ASTNode, error) {
  if err := discardLParen(parsingContext); err != nil { return nil, fmt.Errorf("failed to parse let value: %w", err) }
	if err := discardLet(parsingContext); err != nil { return nil, fmt.Errorf("failed to parse let value: %w", err) }
	if err := discardLParen(parsingContext); err != nil { return nil, fmt.Errorf("failed to parse let value: %w", err) }

	envVariableName, err := parseSymbol(parsingContext)
	if err != nil { return nil, fmt.Errorf("failed to parse let value, while parsing env variable name: %w", err) }

	envVariableValue, err := parseSingle(parsingContext)
  if err != nil { return nil, fmt.Errorf("failed to parse let value, while parsing env variable value: %w", err) }
  if err := discardRParen(parsingContext); err != nil { return nil, fmt.Errorf("failed to parse let value: %w", err) }

	if err := discardIN(parsingContext); err != nil { return nil, fmt.Errorf("failed to parse let value, try to discard 'in': %w", err) } // discardIN(parsingContext)

	body, err := parseSingle(parsingContext)
  if err != nil { return nil, fmt.Errorf("failed to parse let value, while parsing body: %w", err) }
  if err := discardRParen(parsingContext); err != nil { return nil, fmt.Errorf("failed to parse let value, try consume last rparen: %w", err) }

	return &LetNode{
		LetEnv: map[string]ASTNode{envVariableName.Name: envVariableValue},
		Body: body,
	}, nil
}

func discardLet(parsingContext *ParsingContext) error {
  token := parsingContext.consume()
  if token.Type != LET {
    return fmt.Errorf("expected let but get:  %v", token)
  }

  return nil
}

func discardIN(parsingContext *ParsingContext) error {
  token := parsingContext.consume()
  if token.Type != IN {
    return fmt.Errorf("expected in but get:  %v", token)
  }

  return nil
}
