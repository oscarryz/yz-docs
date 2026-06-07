package main

import std "yz/runtime/rt"

type _OptionVariant int

const (
	_OptionSome _OptionVariant = iota
	_OptionNone
)

type Option struct {
	_variant _OptionVariant
	v *Node
}

func NewOptionSome(v *Node) *Option {
	return &Option{
		_variant: _OptionSome,
		v: v,
	}
}

func NewOptionNone() *Option {
	return &Option{
		_variant: _OptionNone,
	}
}

func (self *Option) String() string {
	switch self._variant {
	case _OptionSome:
		return "Option.Some(v: " + std.StringifyRepr(self.v) + ")"
	case _OptionNone:
		return "Option.None()"
	}
	return "Option(?)"
}

type Node struct {
	std.Cown
	value std.Int
	next *Option
}

func NewNode(value std.Int, next *Option) *Node {
	return &Node{
		value: value,
		next: next,
	}
}

func (self *Node) String() string {
	return "Node(value: " + std.StringifyRepr(self.value) + ", next: " + std.StringifyRepr(self.next) + ")"
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
	var a *Node = NewNode(std.NewInt(1), NewOptionNone())
	var b *Node = NewNode(std.NewInt(2), NewOptionNone())
	a.next = NewOptionSome(b)
	std.Print(std.NewString(std.StringifyRepr(a)))
	return std.TheUnit
}

func (self *_mainBoc) Call() std.Unit {
	return std.LazyUnit(std.Schedule(&self.Cown, func() std.Unit {
		return self.call()
	}))
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
