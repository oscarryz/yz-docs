package main

import std "yz/runtime/rt"

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.Schedule(&self.Cown, func() std.Unit {
		var x std.Int = std.NewInt(0)
		if x.Gt(std.NewInt(0)).GoBool() {
			std.Print(std.NewString("positive"))
		} else if x.Lt(std.NewInt(0)).GoBool() {
			std.Print(std.NewString("negative"))
		} else {
			std.Print(std.NewString("zero"))
		}
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
