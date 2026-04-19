package main

import std "yz/runtime/yzrt"

type _mainBoc struct {
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.Go(func() std.Unit {
		var x std.Int = std.NewInt(5)
		var neg_x std.Int = x.Neg()
		std.Print(neg_x)
		var a std.Int = std.NewInt(10)
		var b std.Int = std.NewInt(3)
		var result std.Int = a.Minus(b.Neg())
		std.Print(result)
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
