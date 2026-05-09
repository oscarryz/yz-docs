package main

import std "yz/runtime/rt"

func greet(name std.String) *std.Thunk[std.Unit] {
	return std.Go(func() std.Unit {
		return std.Print(name)
	})
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.NewThunk(func() std.Unit {
		_bg0 := &std.BocGroup{}
		std.Schedule(&self.Cown, func() std.Unit {
			_bg0.Go(func() any {
				return greet(std.NewString("Alice")).Force()
			})
			_bg0.Go(func() any {
				return greet(std.NewString("Bob")).Force()
			})
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
