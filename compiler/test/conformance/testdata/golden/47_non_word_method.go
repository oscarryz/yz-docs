package main

import std "yz/runtime/rt"

type Greeter struct {
	std.Cown
	name std.String
}

func NewGreeter(name std.String) *Greeter {
	return &Greeter{
		name: name,
	}
}

func (self *Greeter) Plusplus(other std.String) *std.Thunk[std.String] {
	return std.Schedule(&self.Cown, func() std.String {
		return std.NewString(std.Stringify(self.name)).Plus(std.NewString(" and ")).Plus(std.NewString(std.Stringify(other)))
	})
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.Schedule(&self.Cown, func() std.Unit {
		var a *Greeter = NewGreeter(std.NewString("Ann"))
		c := a.Plusplus(std.NewString("Taylor"))
		std.Print(c.Force())
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
