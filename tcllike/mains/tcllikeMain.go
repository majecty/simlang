package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"

	"simlang/tcllike/evaluator"
	"simlang/tcllike/lexer"
	"simlang/tcllike/parser"
	"simlang/tcllike/ui"
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

		// Lexer 과정 출력
		tokens := lexer.Tokenize(input)
		fmt.Println("Tokens:")
		for _, token := range tokens {
			fmt.Printf("  %s: %q\n", token.Type, token.Value)
		}

		// Parser 과정 출력
		ast, err := parser.Parse(tokens)
		if err != nil {
			ui.PrintError(err.Error())
			continue
		}
		fmt.Println("AST:")
		fmt.Println(ast.String())

		// Eval 과정
		result, err := evaluator.Eval(ast)
		if err != nil {
			ui.PrintError(err.Error())
			continue
		}

		ui.PrintResult(fmt.Sprintf("%v", result))
	}
}

func runWebUI() {
	fmt.Println("Starting Tcl-like web server on http://localhost:8080")
	http.Handle("/", ui.NewWebUI())
	http.ListenAndServe(":8080", nil)
}
