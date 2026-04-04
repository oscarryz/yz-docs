package main

import std "yz/runtime/yzrt"

type _OptionVariant int

const (
	_OptionSome _OptionVariant = iota
	_OptionNone
)

type Option[V any] struct {
	_variant _OptionVariant
	value V
}

func NewOptionSome[V any](value V) *Option[V] {
	return &Option[V]{
		_variant: _OptionSome,
		value: value,
	}
}

func NewOptionNone[V any]() *Option[V] {
	return &Option[V]{
		_variant: _OptionNone,
	}
}


func main() {
	x := NewOptionSome(std.NewString("hello"))
	switch x._variant {
	case _OptionSome:
		std.Print(x.value)
	case _OptionNone:
		std.Print(std.NewString("nothing"))
	}
}
