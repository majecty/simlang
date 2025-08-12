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
