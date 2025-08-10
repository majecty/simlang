package main

import "fmt"

type IRValue interface {
	toIRValue() string
}

type NumberLiteral struct {
	Value float64
}

func (n *NumberLiteral) toIRValue() string {
	return fmt.Sprintf("%f", n.Value)
}

type RegisterName struct {
	Name string
}

func (r *RegisterName) toIRValue() string {
	return r.Name
}
