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

func (self *Greeter) String() string {
	return "Greeter(name: " + std.Stringify(self.name) + ")"
}

func (self *Greeter) plusplus(other std.String) std.String {
	return std.NewString(std.Stringify(self.name)).Plus(std.NewString(" and ")).Plus(std.NewString(std.Stringify(other)))
}

func (self *Greeter) Plusplus(other std.String) *std.Thunk[std.String] {
	return std.Schedule(&self.Cown, func() std.String {
		return self.plusplus(other)
	})
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.NewThunk(func() std.Unit {
		_bg0 := &std.BocGroup{}
		var c any
		var a *Greeter
		std.Schedule(&self.Cown, func() std.Unit {
			a = NewGreeter(std.NewString("Ann"))
			std.GoStore(_bg0, a.Plusplus(std.NewString("Taylor")), &c)
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		std.Print(c)
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
