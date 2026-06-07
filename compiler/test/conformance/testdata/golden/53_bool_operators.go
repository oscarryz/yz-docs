package main

import std "yz/runtime/rt"

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
	var a std.Bool = std.NewBool(true)
	var b std.Bool = std.NewBool(false)
	std.Print(std.NewString(std.StringifyRepr(a.Ampamp(b))))
	std.Print(std.NewString(std.StringifyRepr(a.Pipepipe(b))))
	std.Print(std.NewString(std.StringifyRepr(a.Ampamp(a))))
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
