package main

import std "yz/runtime/rt"

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) foo() std.String {
	return std.NewString("hello")
}

func (self *_mainBoc) Foo() std.String {
	return std.LazyString(std.Schedule(&self.Cown, func() std.String {
		return self.foo()
	}))
}

func (self *_mainBoc) bar() std.String {
	return std.NewString("world")
}

func (self *_mainBoc) Bar() std.String {
	return std.LazyString(std.Schedule(&self.Cown, func() std.String {
		return self.bar()
	}))
}

func (self *_mainBoc) call() std.Unit {
	std.Print(self.Foo())
	std.Print(self.Bar())
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
