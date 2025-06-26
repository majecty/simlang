package main

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
		case NUMBER:
			current = append(current, token.Value)
		}
	}

	ast = current
	return ast
}
