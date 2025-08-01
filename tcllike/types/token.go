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
	LPAREN
	RPAREN
)

func (t TokenType) String() string {
	switch t {
	case ATOM:
		return "ATOM"
	case NUMBER:
		return "NUMBER"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	default:
		return "UNKNOWN"
	}
}
