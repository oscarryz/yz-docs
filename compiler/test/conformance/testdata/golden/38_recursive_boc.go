package main

import std "yz/runtime/yzrt"

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

func main() {
	_bg0 := &std.BocGroup{}
	_bg0.Go(func() any {
		return countdown(std.NewInt(3)).Force()
	})
	_bg0.Wait()
}
