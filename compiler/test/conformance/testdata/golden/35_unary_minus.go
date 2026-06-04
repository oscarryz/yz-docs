package main

import std "yz/runtime/rt"

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
	var x std.Int = std.NewInt(5)
	var neg_x std.Int = x.Neg()
	std.Print(std.NewString(std.StringifyRepr(neg_x)))
	var a std.Int = std.NewInt(10)
	var b std.Int = std.NewInt(3)
	var result std.Int = a.Minus(b.Neg())
	std.Print(std.NewString(std.StringifyRepr(result)))
	return std.TheUnit
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.Schedule(&self.Cown, func() std.Unit {
		return self.call()
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
