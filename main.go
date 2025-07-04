package main

import "fmt"

func main() {
	fmt.Println("Hello, Go Project!")

	fmt.Println(parse(tokenize("(hello world)")))

	// 덧셈 테스트
	if result, err := parse(tokenize("(+ 1 2)")); err != nil {
    panic(err)
  } else {
    fmt.Println("(+ 1 2) =", eval(result))
	}

	if result, err := parse(tokenize("(+ 1 2 3)")); err != nil {
		panic(err)
	} else  {
    fmt.Println("(+ 1 2 3) =", eval(result))
	}

	if rsult, err := parse(tokenize("(+ 10 (+ 5 3) 2)")); err != nil {
    panic(err)
  } else {
    fmt.Println("(+ 10 (+ 5 3) 2) =", eval(rsult))
  }

	if result, err := parse(tokenize("(+ -5 10)")); err != nil {
    panic(err)
  } else {
    fmt.Println("(+ -5 10) =", eval(result))
  }

  if result, err := parse(tokenize("(+ 0 0 0)")); err != nil {
    panic(err)
  } else {
    fmt.Println("(+ 0 0 0) =", eval(result))
  }

	if result, err := parse(tokenize("(+ 100 -50 25)")); err != nil {
    panic(err)
  } else {
    fmt.Println("(+ 100 -50 25) =", eval(result))
  }

	if result, err := parse(tokenize("(+ (+ 1 2) (+ 3 4) (+ 5 6))")); err != nil {
    panic(err)
	} else {
    fmt.Println("(+ (+ 1 2) (+ 3 4) (+ 5 6)) =", eval(result))
  }

	if result, err := parse(tokenize("(+ 42)")); err != nil {
    panic(err)
  } else {
    fmt.Println("(+ 42) =", eval(result))
  }

	if result, err := parse(tokenize("(let (x 10) in x)")); err != nil {
    panic(err)
  } else {
    fmt.Println("(let (x 10) in x) =", eval(result))
  }

  if result, err := parse(tokenize("(let (x 10) in (let (y 20) in (+ x y)))")); err != nil {
    panic(err)
  } else {
    fmt.Println("(let (x 10) in (let (y 20) in (+ x y))) =", eval(result))
  }

	if result, err := parse(tokenize("(let (x 10) in (let (x 20) in x))")); err != nil {
    panic(err)
  } else {
    fmt.Println("(let (x 10) in (let (x 20) in x)) =", eval(result))
  }
}
