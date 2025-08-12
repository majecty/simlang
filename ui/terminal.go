package ui

import (
	"fmt"
	"strings"
	"github.com/charmbracelet/lipgloss"
)

var (
	promptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("36")).Bold(true)
	resultStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	errorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
)

func PrintPrompt() {
	fmt.Print(promptStyle.Render("λ > "))
}

func PrintResult(result string) {
	fmt.Println(resultStyle.Render("=> " + result))
}

func PrintError(err string) {
	fmt.Println(errorStyle.Render("✗ " + err))
}

func PrintWelcome() {
	welcome := `
  ____  _       _   _     
 / ___|(_)_ __ | |_| |__  
 \___ \| | '_ \| __| '_ \ 
  ___) | | | | | |_| | | |
 |____/|_|_| |_|\__|_| |_|

`
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Render(welcome))
	fmt.Println("간단한 Lisp 인터프리터 (종료하려면 exit 입력)")
	fmt.Println(strings.Repeat("=", 40))
}
