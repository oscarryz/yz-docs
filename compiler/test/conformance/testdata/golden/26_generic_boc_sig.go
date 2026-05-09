package main

import std "yz/runtime/rt"

func identity[V any](value V) *std.Thunk[V] {
	return std.Go(func() V {
		return value
	})
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.Schedule(&self.Cown, func() std.Unit {
		x := identity(std.NewString("hello"))
		std.Print(x.Force())
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
