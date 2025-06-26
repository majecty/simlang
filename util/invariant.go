package util

import "fmt"

func Invariant(predicate bool, message string, args ...any) {
	if !predicate {
		panic(fmt.Sprintf(message, args...))
	}
}
