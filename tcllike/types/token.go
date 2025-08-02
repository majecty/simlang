// Package types implements types used in the Simlang programming language.
package types

type Token struct {
	Type  TokenType
	Value string
}

type TokenType int

const (
	Atom TokenType = iota
	Number
	LParen
	RParen
	LineEnd
	LBracket
	RBracket
)

func (t TokenType) String() string {
	switch t {
	case Atom:
		return "ATOM"
	case Number:
		return "NUMBER"
	case LParen:
		return "LPAREN"
	case RParen:
		return "RPAREN"
	case LineEnd:
		return "LineEnd"
	case LBracket:
		return "LBracket"
	case RBracket:
		return "RBracket"
	default:
		return "UNKNOWN"
	}
}
