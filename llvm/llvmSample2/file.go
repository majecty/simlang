package main

import (
	"errors"
	"fmt"
	"os"
)

func writeToFile(filename string, content string) (err error) {
	// Create or open the file for writing
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", filename, err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			errors.Join(err, fmt.Errorf("error closing file %s: %v", filename, closeErr))
		}
	}()

	// Write the LLVM IR string to the file
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write LLVM IR to file %s: %v", filename, err)
	}
	return nil
}
