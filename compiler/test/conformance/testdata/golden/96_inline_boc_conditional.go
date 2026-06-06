package main

import std "yz/runtime/rt"

type _greeterBoc struct {
	std.Cown
}

func (self *_greeterBoc) String() string {
	return "{ " + "greet: {}" + " }"
}

func (self *_greeterBoc) greet() std.String {
	return std.NewString("hello")
}

func (self *_greeterBoc) Greet() *std.Thunk[std.String] {
	return std.Schedule(&self.Cown, func() std.String {
		return self.greet()
	})
}

var Greeter = &_greeterBoc{}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.NewThunk(func() std.Unit {
		_bg0 := &std.BocGroup{}
		std.Schedule(&self.Cown, func() std.Unit {
			_st0 := std.WrapStringThunk(Greeter.Greet()).Eqeq(std.NewString("hello")).Qm(func() any {
				std.Print(std.NewString("matched"))
				return std.TheUnit
			}, func() any {
				std.Print(std.NewString("no match"))
				return std.TheUnit
			})
			_bg0.Add(func() { _st0.Force() })
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
