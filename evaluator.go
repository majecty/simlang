package main

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
