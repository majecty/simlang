// Package types implements types used in the Simlang programming language.
package types

type Token struct {
	Type  TokenType
	Value string
}

type TokenType int

const (
	ATOM TokenType = iota
	NUMBER
)

func (t TokenType) String() string {
	switch t {
	case ATOM:
		return "ATOM"
	case NUMBER:
		return "NUMBER"
	default:
		return "UNKNOWN"
	}
}
