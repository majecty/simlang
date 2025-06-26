package main

import "fmt"

func main() {
	fmt.Println("Hello, Go Project!")

	fmt.Println(parse(tokenize("(hello world)")))
}

type TokenType int

const (
	LPAREN TokenType = iota
	RPAREN
	ATOM
)

func (t TokenType) String() string {
	switch t {
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case ATOM:
		return "ATOM"
	default:
		return "UNKNOWN"
	}
}

type Token struct {
	Type  TokenType
	Value string
}

func tokenize(input string) []Token {
	tokens := []Token{}
	var current string

	for i := 0; i < len(input); i++ {
		ch := input[i]

		switch ch {
		case '(':
			tokens = append(tokens, Token{LPAREN, "("})
		case ')':
			if current != "" {
				tokens = append(tokens, Token{ATOM, current})
				current = ""
			}
			tokens = append(tokens, Token{RPAREN, ")"})
		case ' ', '\n', '\t':
			if current != "" {
				tokens = append(tokens, Token{ATOM, current})
				current = ""
			}
		default:
			current += string(ch)
		}
	}

	return tokens
}

func parse(tokens []Token) []any {
	var ast []any
	var stack [][]any
	current := []any{}

	for _, token := range tokens {
		switch token.Type {
		case LPAREN:
			stack = append(stack, current)
			current = []any{}
		case RPAREN:
			if len(stack) > 0 {
				last := current
				current = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				current = append(current, last)
			}
		case ATOM:
			current = append(current, token.Value)
		}
	}

	ast = current
	return ast
}
