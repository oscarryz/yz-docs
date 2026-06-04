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
	var r1 *std.Option[std.Int] = d.AtOpt(std.NewString("a"))
	switch r1.Variant {
	case std.OptionSome:
		std.Print(std.NewString(std.StringifyRepr(r1.Value)))
	case std.OptionNone:
		std.Print(std.NewString("not found"))
	}
	var r2 *std.Option[std.Int] = d.AtOpt(std.NewString("z"))
	if std.NewBool(r2.Variant == std.OptionSome).GoBool() {
		std.Print(std.NewString(std.StringifyRepr(r2.Value)))
	}
	if std.NewBool(r2.Variant == std.OptionNone).GoBool() {
		std.Print(std.NewString("absent"))
	}
	std.Print(d.AtOpt(std.NewString("b")).ToStr())
	var r3 *std.Option[std.Int] = d.AtOpt(std.NewString("a"))
	std.Print(std.NewString(std.StringifyRepr(r3)))
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
