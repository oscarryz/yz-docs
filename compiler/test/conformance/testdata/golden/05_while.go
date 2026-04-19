package main

import std "yz/runtime/yzrt"

type _mainBoc struct {
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.Go(func() std.Unit {
		var n std.Int = std.NewInt(0)
		for n.Lt(std.NewInt(3)).GoBool() {
			n = n.Plus(std.NewInt(1))
		}
		std.Print(n)
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
