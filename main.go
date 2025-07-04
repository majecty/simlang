package main

import "fmt"

func main() {
	fmt.Println("Hello, Go Project!")

	fmt.Println(parse(tokenize("(hello world)")))

	// 덧셈 테스트
	if result, err := parse(tokenize("(+ 1 2)")); err != nil {
    panic(err)
  } else if evalResult, err := eval(result); err != nil {
    panic(err)
  } else {
    fmt.Println("(+ 1 2) =", evalResult)
	}

	if result, err := parse(tokenize("(+ 1 2 3)")); err != nil {
		panic(err)
	} else if evalResult, err := eval(result); err != nil {
    panic(err)
  } else {
    fmt.Println("(+ 1 2 3) =", evalResult)
	}

	if rsult, err := parse(tokenize("(+ 10 (+ 5 3) 2)")); err != nil {
    panic(err)
	} else if evalResult, err := eval(rsult); err != nil {
    panic(err)
  } else {
    fmt.Println("(+ 10 (+ 5 3) 2) =", evalResult)
  }

	if result, err := parse(tokenize("(+ -5 10)")); err != nil {
    panic(err)
	} else if evalResult, err := eval(result); err != nil {
    panic(err)
  } else {
    fmt.Println("(+ -5 10) =", evalResult)
  }

  if result, err := parse(tokenize("(+ 0 0 0)")); err != nil {
    panic(err)
	} else if evalResult, err := eval(result); err != nil {
    panic(err)
  } else {
    fmt.Println("(+ 0 0 0) =", evalResult)
  }

	if result, err := parse(tokenize("(+ 100 -50 25)")); err != nil {
    panic(err)
	} else if evalResult, err := eval(result); err != nil {
    panic(err)
  } else {
    fmt.Println("(+ 100 -50 25) =", evalResult)
  }

	if result, err := parse(tokenize("(+ (+ 1 2) (+ 3 4) (+ 5 6))")); err != nil {
    panic(err)
	} else if evalResult, err := eval(result); err != nil {
    panic(err)
	} else {
    fmt.Println("(+ (+ 1 2) (+ 3 4) (+ 5 6)) =", evalResult)
  }

	if result, err := parse(tokenize("(+ 42)")); err != nil {
    panic(err)
	} else if evalResult, err := eval(result); err != nil {
    panic(err)
  } else {
    fmt.Println("(+ 42) =", evalResult)
  }

	if result, err := parse(tokenize("(let (x 10) in x)")); err != nil {
    panic(err)
	} else if evalResult, err := eval(result); err != nil {
    panic(err)
  } else {
    fmt.Println("(let (x 10) in x) =", evalResult)
  }

  if result, err := parse(tokenize("(let (x 10) in (let (y 20) in (+ x y)))")); err != nil {
    panic(err)
	} else if evalResult, err := eval(result); err != nil {
    panic(err)
  } else {
    fmt.Println("(let (x 10) in (let (y 20) in (+ x y))) =", evalResult)
  }

	if result, err := parse(tokenize("(let (x 10) in (let (x 20) in x))")); err != nil {
    panic(err)
	} else if evalResult, err := eval(result); err != nil {
    panic(err)
  } else {
    fmt.Println("(let (x 10) in (let (x 20) in x)) =", evalResult)
  }

	if result, err := parse(tokenize("(let (add2 (lambda (x) (+ x 2))) in (add2 10))")); err != nil {
    panic(err)
	} else if evalResult, err := eval(result); err != nil {
    panic(err)
  } else {
    fmt.Println("(let (add2 (lambda (x) (+ x 2))) in (add2 10)) =", evalResult)
  }
}
