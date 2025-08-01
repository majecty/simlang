// Package lexer implements lexical analysis for the Simlang programming language.
package lexer

import "simlang/tcllike/types"

func Tokenize(input string) []types.Token {
	tokens := []types.Token{}
	var current string

	for i := 0; i < len(input); i++ {
		ch := input[i]

		switch ch {
		case ' ', '\n', '\t':
			if current != "" {
				tokens = append(tokens, createToken(current))
				current = ""
			}
		default:
			current += string(ch)
		}
	}

	if current != "" {
		tokens = append(tokens, createToken(current))
	}

	return tokens
}

func createToken(value string) types.Token {
	if isNumber(value) {
		return types.Token{Type: types.NUMBER, Value: value}
	}
	return types.Token{Type: types.ATOM, Value: value}
}

func isNumber(s string) bool {
	if len(s) == 0 {
		return false
	}

	start := 0
	if s[0] == '-' || s[0] == '+' {
		if len(s) == 1 {
			return false
		}
		start = 1
	}

	for i := start; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}
