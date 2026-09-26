package main

import std "yz/runtime/rt"

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
	var d std.Dict[std.String, std.Int] = std.NewDict[std.String, std.Int]().Set(std.NewString("a"), std.NewInt(1)).Set(std.NewString("b"), std.NewInt(2)).Set(std.NewString("c"), std.NewInt(3))
	d.Each(func(k std.String, v std.Int) std.Unit {
		return std.Print(k.ToStr().Plus(std.NewString(" = ")).Plus(v.ToStr()))
	})
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
