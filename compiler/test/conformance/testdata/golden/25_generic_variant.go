package main

import std "yz/runtime/rt"

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


type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.Schedule(&self.Cown, func() std.Unit {
		x := NewOptionSome(std.NewString("hello"))
		switch x._variant {
		case _OptionSome:
			std.Print(x.value)
		case _OptionNone:
			std.Print(std.NewString("nothing"))
		}
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
