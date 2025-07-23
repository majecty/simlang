package lexer

import (
	"simlang/types"
)


func Toknize(input string) []types.Token {
	tokens := []types.Token{}
	var current string

	for i := 0; i < len(input); i++ {
		ch := input[i]

		switch ch {
		case '(':
			tokens = append(tokens, types.Token{Type: types.LPAREN, Value: "("})
		case ')':
			if current != "" {
				tokens = append(tokens, createToken(current))
				current = ""
			}
			tokens = append(tokens, types.Token{Type: types.RPAREN, Value:  ")"})
		case ' ', '\n', '\t':
			if current != "" {
				tokens = append(tokens, createToken(current))
				current = ""
			}
		default:
			current += string(ch)
		}
	}

	return tokens
}

func createToken(value string) types.Token {
	// 숫자인지 확인
	if isNumber(value) {
		return types.Token{Type: types.NUMBER, Value: value}
	}

	switch value {
	case "let":
		return types.Token{Type: types.LET, Value: value}
	case "in":
		return types.Token{Type: types.IN, Value:  value}
	case "lambda":
		return types.Token{Type: types.LAMBDA, Value:  value}
	}

	return types.Token{Type: types.ATOM, Value:  value}
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
