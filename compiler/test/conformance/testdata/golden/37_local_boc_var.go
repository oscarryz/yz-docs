package main

import std "yz/runtime/yzrt"

type _mainBoc struct {
}

func (self *_mainBoc) Foo() *std.Thunk[std.String] {
	return std.Go(func() std.String {
		return std.NewString("hello")
	})
}

func (self *_mainBoc) Bar() *std.Thunk[std.String] {
	return std.Go(func() std.String {
		return std.NewString("world")
	})
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.Go(func() std.Unit {
		std.Print(self.Foo().Force())
		std.Print(self.Bar().Force())
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
