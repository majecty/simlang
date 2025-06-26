package main

type TokenType int

const (
	LPAREN TokenType = iota
	RPAREN
	ATOM
	NUMBER
)

func (t TokenType) String() string {
	switch t {
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case ATOM:
		return "ATOM"
	case NUMBER:
		return "NUMBER"
	default:
		return "UNKNOWN"
	}
}

type Token struct {
	Type  TokenType
	Value string
}
