package main

import std "yz/runtime/rt"

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
	var x std.Int = std.NewInt(15)
	if x.Gt(std.NewInt(0)).GoBool() {
		std.Print(std.NewString("positive"))
		std.Print(std.NewString("multi-statement arm"))
	} else if x.Gt(std.NewInt(10)).GoBool() {
		std.Print(std.NewString("and greater than 10"))
	} else {
		std.Print(std.NewString("done"))
	}
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
