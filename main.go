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
