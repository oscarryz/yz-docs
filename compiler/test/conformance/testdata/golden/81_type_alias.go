package main

import std "yz/runtime/rt"

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

type Bar = Foo


type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
	var b1 *Foo = NewFoo(std.NewString("Alice"))
	var b2 *Foo = NewFoo(std.NewString("Bob"))
	std.Print(b1.name)
	std.Print(b2.name)
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
