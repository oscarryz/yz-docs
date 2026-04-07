package main

import std "yz/runtime/yzrt"

type Box[T any] struct {
	value T
}

func NewBox[T any](value T) *Box[T] {
	return &Box[T]{
		value: value,
	}
}

func main() {
	b := NewBox(std.NewInt(42))
	std.Print(b.value)
}
