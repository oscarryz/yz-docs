package main

import std "yz/runtime/rt"

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) call() std.Unit {
	var x std.Int = std.NewInt(5)
	var neg_x std.Int = x.Neg()
	std.Print(neg_x)
	var a std.Int = std.NewInt(10)
	var b std.Int = std.NewInt(3)
	var result std.Int = a.Minus(b.Neg())
	std.Print(result)
	return std.TheUnit
}

func (self *_mainBoc) Call() std.Unit {
	return std.LazyUnit(std.Schedule(&self.Cown, func() std.Unit {
		return self.call()
	}))
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
