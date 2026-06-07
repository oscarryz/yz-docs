package main

import std "yz/runtime/rt"

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
	var d std.Dict[std.String, std.Int] = std.NewDict[std.String, std.Int]().Set(std.NewString("a"), std.NewInt(1)).Set(std.NewString("b"), std.NewInt(2))
	d = d.Set(std.NewString("c"), std.NewInt(3))
	var c *std.Option[std.Int] = d.AtOpt(std.NewString("c"))
	if std.NewBool(c.Variant == std.OptionSome).GoBool() {
		std.Print(std.NewString(std.StringifyRepr(c.Value)))
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
