package main

import std "yz/runtime/rt"

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) call() std.Unit {
	var x std.Int = std.NewInt(5)
	if x.Gt(std.NewInt(3)).GoBool() {
		std.Print(std.NewString("big"))
	} else {
		std.Print(std.NewString("small"))
	}
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
