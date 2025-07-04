package main

type Token struct {
	Type  TokenType
	Value string
}

type TokenType int

const (
	LPAREN TokenType = iota
	RPAREN
	ATOM
	NUMBER
	LET
	IN // let in
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
	case LET:
    return "LET"
	case IN:
    return "IN"
	default:
		return "UNKNOWN"
	}
}
