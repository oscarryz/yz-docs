package main

import std "yz/runtime/rt"

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) call() std.Unit {
	var name std.String = std.NewString("World")
	var n std.Int = std.NewInt(42)
	std.Print(std.NewString("Hello, ").Plus(std.NewString(std.Stringify(name))).Plus(std.NewString("!")))
	std.Print(std.NewString("Answer: ").Plus(std.NewString(std.Stringify(n))))
	std.Print(std.NewString("Sum: ").Plus(std.NewString(std.Stringify(n.Plus(std.NewInt(1))))))
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
