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
		var x V
		_bgs_x := &std.BocGroup{}
		std.Schedule(&self.Cown, func() std.Unit {
			std.GoStore(_bgs_x, identity(std.NewString("hello")), &x)
			return std.TheUnit
		}).Force()
		_bgs_x.Wait()
		std.Print(x)
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
