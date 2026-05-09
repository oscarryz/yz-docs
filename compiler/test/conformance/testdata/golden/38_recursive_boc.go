package main

import std "yz/runtime/rt"

func countdown(n std.Int) *std.Thunk[std.Unit] {
	return std.Go(func() std.Unit {
		if n.Eqeq(std.NewInt(0)).GoBool() {
			std.Print(std.NewString("done"))
		} else {
			std.Print(n)
			countdown(n.Minus(std.NewInt(1))).Force()
		}
		return std.TheUnit
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
				return countdown(std.NewInt(3)).Force()
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
