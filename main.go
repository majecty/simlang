package main

import "fmt"

func main() {
	fmt.Println("Hello, Go Project!")

	fmt.Println(parse(tokenize("(hello world)")))

	// 덧셈 테스트
	result := eval(parse(tokenize("(+ 1 2 3)")))
	fmt.Println("(+ 1 2 3) =", result)

	result2 := eval(parse(tokenize("(+ 10 (+ 5 3) 2)")))
	fmt.Println("(+ 10 (+ 5 3) 2) =", result2)

	// 더 많은 덧셈 테스트 케이스
	result3 := eval(parse(tokenize("(+ -5 10)")))
	fmt.Println("(+ -5 10) =", result3)

	result4 := eval(parse(tokenize("(+ 0 0 0)")))
	fmt.Println("(+ 0 0 0) =", result4)

	result5 := eval(parse(tokenize("(+ 100 -50 25)")))
	fmt.Println("(+ 100 -50 25) =", result5)

	result6 := eval(parse(tokenize("(+ (+ 1 2) (+ 3 4) (+ 5 6))")))
	fmt.Println("(+ (+ 1 2) (+ 3 4) (+ 5 6)) =", result6)

	result7 := eval(parse(tokenize("(+ 42)")))
	fmt.Println("(+ 42) =", result7)
}

type TokenType int

const (
	LPAREN TokenType = iota
	RPAREN
	ATOM
	NUMBER
)

func (t TokenType) String() string {
	switch t {
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case ATOM:
		return "ATOM"
	case NUMBER:
		return "NUMBER"
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

func eval(ast []any) any {
	if len(ast) == 0 {
		return nil
	}

	// 단일 요소인 경우
	if len(ast) == 1 {
		return evalSingle(ast[0])
	}

	// 리스트인 경우 첫 번째 요소가 함수/연산자
	first := ast[0]

	switch first {
	case "+":
		return evalAdd(ast[1:])
	default:
		// 다른 함수들은 나중에 구현
		return ast
	}
}

func evalSingle(item any) any {
	switch v := item.(type) {
	case string:
		// 숫자 문자열을 정수로 변환
		if isNumber(v) {
			if num := parseInt(v); num != nil {
				return *num
			}
		}
		return v
	case []any:
		// 중첩된 리스트를 재귀적으로 평가
		return eval(v)
	default:
		return v
	}
}

func evalAdd(args []any) int {
	sum := 0
	for _, arg := range args {
		evaluated := evalSingle(arg)
		if num, ok := evaluated.(int); ok {
			sum += num
		}
	}
	return sum
}

func parseInt(s string) *int {
	if !isNumber(s) {
		return nil
	}

	result := 0
	negative := false
	start := 0

	if s[0] == '-' {
		negative = true
		start = 1
	} else if s[0] == '+' {
		start = 1
	}

	for i := start; i < len(s); i++ {
		digit := int(s[i] - '0')
		result = result*10 + digit
	}

	if negative {
		result = -result
	}

	return &result
}
