package main

import std "yz/runtime/rt"

type Bar struct {
	std.Cown
	f std.String
}

func NewBar(f std.String) *Bar {
	return &Bar{
		f: f,
	}
}

func (self *Bar) String() string {
	return "Bar(f: " + std.StringifyRepr(self.f) + ")"
}

func (self *Bar) F() std.String {
	return self.f
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
	var b *Bar = &Bar{}
	b.f = std.NewString("hello")
	std.Print(b.f)
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
