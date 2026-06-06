package main

import std "yz/runtime/rt"

type _greetingBoc struct {
	std.Cown
}

func (self *_greetingBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_greetingBoc) call() std.String {
	return std.NewString("Hello from singleton")
}

func (self *_greetingBoc) Call() *std.Thunk[std.String] {
	return std.Schedule(&self.Cown, func() std.String {
		return self.call()
	})
}

var Greeting = &_greetingBoc{}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.NewThunk(func() std.Unit {
		_bg0 := &std.BocGroup{}
		var msg std.String
		std.Schedule(&self.Cown, func() std.Unit {
			_st0 := Greeting.Call()
			_bg0.Add(func() { msg = _st0.Force() })
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		std.Print(msg)
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
