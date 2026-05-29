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

func (self *Option[V]) String() string {
	switch self._variant {
	case _OptionSome:
		return "Option.Some(value: " + std.StringifyRepr(self.value) + ")"
	case _OptionNone:
		return "Option.None()"
	}
	return "Option(?)"
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
	var x *Option[std.String] = NewOptionSome(std.NewString("hello"))
	switch x._variant {
	case _OptionSome:
		std.Print(x.value)
	case _OptionNone:
		std.Print(std.NewString("nothing"))
	}
	return std.TheUnit
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.Schedule(&self.Cown, func() std.Unit {
		return self.call()
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
