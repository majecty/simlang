package main

import (
	"fmt"
	"log"
	"os"
)

// https://llvm.org/docs/LangRef.html
func main() {
	fmt.Println("Hello, Go Project!")

	llvmIR := `; ModuleID = 'simple_module'
source_filename = "simple_program.ll"

define i32 @main() {
  ret i32 42
}`

	filename := "output.ll"

	// Create or open the file for writing
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to create file %s: %v", filename, err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Printf("Error closing file %s: %v", filename, closeErr)
		}
	}()

	// Write the LLVM IR string to the file
	_, err = file.WriteString(llvmIR)
	if err != nil {
		log.Fatalf("Failed to write LLVM IR to file %s: %v", filename, err)
	}

	log.Printf("Successfully wrote LLVM IR to %s", filename)

}
