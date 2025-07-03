package main

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
				tokens = append(tokens, createToken(current))
				current = ""
			}
			tokens = append(tokens, Token{RPAREN, ")"})
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

func createToken(value string) Token {
	// 숫자인지 확인
	if isNumber(value) {
		return Token{NUMBER, value}
	}

	switch value {
	case "let":
    return Token{LET, value}
	}

	return Token{ATOM, value}
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
