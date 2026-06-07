package main

import std "yz/runtime/rt"

type _greetBoc struct {
	std.Cown
	name std.String
}

func (self *_greetBoc) String() string {
	return "{ " + "name: " + std.StringifyRepr(self.name) + "; " + "call: {}" + " }"
}

func (self *_greetBoc) Call(name std.String) std.Unit {
	return std.LazyUnit(std.Schedule(&self.Cown, func() std.Unit {
		self.name = name
		return std.Print(self.name)
	}))
}

var Greet = &_greetBoc{
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) Call() std.Unit {
	return std.LazyUnit(std.NewThunk(func() std.Unit {
		_bg0 := &std.BocGroup{}
		std.Schedule(&self.Cown, func() std.Unit {
			_st0 := Greet.Call(std.NewString("Alice"))
			_bg0.Add(func() { _st0.Await() })
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		std.Print(Greet.name)
		return std.TheUnit
	}))
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
