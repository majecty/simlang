package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"simlang/evaluator"
	"simlang/lexer"
	"simlang/parser"
	"simlang/ui"
)

func main() {
	mode := flag.String("mode", "terminal", "Interface mode (terminal/web)")
	flag.Parse()

	switch *mode {
	case "terminal":
		runTerminalUI()
	case "web":
		runWebUI()
	default:
		fmt.Println("Invalid mode. Use 'terminal' or 'web'")
		os.Exit(1)
	}
}

func runTerminalUI() {
	ui.PrintWelcome()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		ui.PrintPrompt()
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		if input == "exit" {
			break
		}

		ast, err := parser.Parse(lexer.Toknize(input))
		if err != nil {
			ui.PrintError(err.Error())
			continue
		}

		result, err := evaluator.Eval(ast)
		if err != nil {
			ui.PrintError(err.Error())
			continue
		}

		ui.PrintResult(fmt.Sprintf("%v", result))
	}
}

func runWebUI() {
	fmt.Println("Starting web server on http://localhost:8080")
	http.Handle("/", ui.NewWebUI())
	http.ListenAndServe(":8080", nil)
}
	fmt.Println("Hello, Go Project!")

	fmt.Println(parser.Parse(lexer.Toknize("(hello world)")))

	// 덧셈 테스트
	if result, err := parser.Parse(lexer.Toknize("(+ 1 2)")); err != nil {
		panic(err)
	} else if evalResult, err := evaluator.Eval(result); err != nil {
		panic(err)
	} else {
		fmt.Println("(+ 1 2) =", evalResult)
	}

	if result, err := parser.Parse(lexer.Toknize("(+ 1 2 3)")); err != nil {
		panic(err)
	} else if evalResult, err := evaluator.Eval(result); err != nil {
		panic(err)
	} else {
		fmt.Println("(+ 1 2 3) =", evalResult)
	}

	if rsult, err := parser.Parse(lexer.Toknize("(+ 10 (+ 5 3) 2)")); err != nil {
		panic(err)
	} else if evalResult, err := evaluator.Eval(rsult); err != nil {
		panic(err)
	} else {
		fmt.Println("(+ 10 (+ 5 3) 2) =", evalResult)
	}

	if result, err := parser.Parse(lexer.Toknize("(+ -5 10)")); err != nil {
		panic(err)
	} else if evalResult, err := evaluator.Eval(result); err != nil {
		panic(err)
	} else {
		fmt.Println("(+ -5 10) =", evalResult)
	}

	if result, err := parser.Parse(lexer.Toknize("(+ 0 0 0)")); err != nil {
		panic(err)
	} else if evalResult, err := evaluator.Eval(result); err != nil {
		panic(err)
	} else {
		fmt.Println("(+ 0 0 0) =", evalResult)
	}

	if result, err := parser.Parse(lexer.Toknize("(+ 100 -50 25)")); err != nil {
		panic(err)
	} else if evalResult, err := evaluator.Eval(result); err != nil {
		panic(err)
	} else {
		fmt.Println("(+ 100 -50 25) =", evalResult)
	}

	if result, err := parser.Parse(lexer.Toknize("(+ (+ 1 2) (+ 3 4) (+ 5 6))")); err != nil {
		panic(err)
	} else if evalResult, err := evaluator.Eval(result); err != nil {
		panic(err)
	} else {
		fmt.Println("(+ (+ 1 2) (+ 3 4) (+ 5 6)) =", evalResult)
	}

	if result, err := parser.Parse(lexer.Toknize("(+ 42)")); err != nil {
		panic(err)
	} else if evalResult, err := evaluator.Eval(result); err != nil {
		panic(err)
	} else {
		fmt.Println("(+ 42) =", evalResult)
	}

	if result, err := parser.Parse(lexer.Toknize("(let (x 10) in x)")); err != nil {
		panic(err)
	} else if evalResult, err := evaluator.Eval(result); err != nil {
		panic(err)
	} else {
		fmt.Println("(let (x 10) in x) =", evalResult)
	}

	if result, err := parser.Parse(lexer.Toknize("(let (x 10) in (let (y 20) in (+ x y)))")); err != nil {
		panic(err)
	} else if evalResult, err := evaluator.Eval(result); err != nil {
		panic(err)
	} else {
		fmt.Println("(let (x 10) in (let (y 20) in (+ x y))) =", evalResult)
	}

	if result, err := parser.Parse(lexer.Toknize("(let (x 10) in (let (x 20) in x))")); err != nil {
		panic(err)
	} else if evalResult, err := evaluator.Eval(result); err != nil {
		panic(err)
	} else {
		fmt.Println("(let (x 10) in (let (x 20) in x)) =", evalResult)
	}

	if result, err := parser.Parse(lexer.Toknize("(let (add2 (lambda (x) (+ x 2))) in (add2 10))")); err != nil {
		panic(err)
	} else if evalResult, err := evaluator.Eval(result); err != nil {
		panic(err)
	} else {
		fmt.Println("(let (add2 (lambda (x) (+ x 2))) in (add2 10)) =", evalResult)
	}
}
