package main

import std "yz/runtime/rt"

type Named struct {
	std.Cown
	name std.String
	greet func() *std.Thunk[std.Unit]
}

func NewNamed(name std.String, greet func() *std.Thunk[std.Unit]) *Named {
	return &Named{
		name: name,
		greet: greet,
	}
}

func (self *Named) String() string {
	return "Named(name: " + std.StringifyRepr(self.name) + ", greet: " + std.StringifyRepr(self.greet) + ")"
}

func (self *Named) Greet() std.Unit {
	return self.greet()
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
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
