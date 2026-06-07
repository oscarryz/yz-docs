package main

import std "yz/runtime/rt"

type _fooBar struct {
	std.Cown
	_outer *Foo
}

func New_fooBar(_outer *Foo) *_fooBar {
	return &_fooBar{
		_outer: _outer,
	}
}

func (self *_fooBar) String() string {
	return "_fooBar(_outer: " + std.StringifyRepr(self._outer) + ")"
}

func (self *_fooBar) describe() std.String {
	return self._outer.name
}

func (self *_fooBar) Describe() std.String {
	return std.LazyString(std.Schedule(&self.Cown, func() std.String {
		return self.describe()
	}))
}

type Foo struct {
	std.Cown
	name std.String
}

func NewFoo(name std.String) *Foo {
	return &Foo{
		name: name,
	}
}

func (self *Foo) String() string {
	return "Foo(name: " + std.StringifyRepr(self.name) + ")"
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
	var alice *Foo = NewFoo(std.NewString("alice"))
	var bob *Foo = NewFoo(std.NewString("bob"))
	var ab *_fooBar = New_fooBar(alice)
	var bb *_fooBar = New_fooBar(bob)
	std.Print(ab.Describe())
	std.Print(bb.Describe())
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
