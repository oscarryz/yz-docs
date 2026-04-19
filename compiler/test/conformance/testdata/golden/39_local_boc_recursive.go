package main

import std "yz/runtime/yzrt"

type _mainBoc struct {
}

func (self *_mainBoc) F(n std.Int) *std.Thunk[std.Unit] {
	return std.Go(func() std.Unit {
		if n.Eqeq(std.NewInt(0)).GoBool() {
			std.Print(std.NewString("fin"))
		} else {
			std.Print(n)
			self.F(n.Minus(std.NewInt(1))).Force()
		}
		return std.TheUnit
	})
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.Go(func() std.Unit {
		_bg0 := &std.BocGroup{}
		_bg0.Go(func() any {
			return self.F(std.NewInt(3)).Force()
		})
		_bg0.Wait()
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
