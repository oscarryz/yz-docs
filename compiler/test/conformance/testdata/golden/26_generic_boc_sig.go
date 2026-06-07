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

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) Call() std.Unit {
	return std.LazyUnit(std.NewThunk(func() std.Unit {
		_bg0 := &std.BocGroup{}
		var x std.String
		std.Schedule(&self.Cown, func() std.Unit {
			_st0 := identity(std.NewString("hello"))
			_bg0.Add(func() { x = _st0.Force() })
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		std.Print(x)
		return std.TheUnit
	}))
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
