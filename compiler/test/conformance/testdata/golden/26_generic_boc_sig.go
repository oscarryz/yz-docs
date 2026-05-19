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
	return std.NewThunk(func() std.Unit {
		_bg0 := &std.BocGroup{}
		var x V
		std.Schedule(&self.Cown, func() std.Unit {
			std.GoStore(_bg0, identity(std.NewString("hello")), &x)
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		std.Print(x)
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
