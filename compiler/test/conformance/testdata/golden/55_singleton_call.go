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

func (self *_greetingBoc) Call() std.String {
	return std.LazyString(std.Schedule(&self.Cown, func() std.String {
		return self.call()
	}))
}

var Greeting = &_greetingBoc{}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) Call() std.Unit {
	return std.LazyUnit(std.NewThunk(func() std.Unit {
		_bg0 := &std.BocGroup{}
		var msg std.String
		std.Schedule(&self.Cown, func() std.Unit {
			msg = Greeting.Call()
			_bg0.Add(func() { msg.Await() })
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		std.Print(msg)
		return std.TheUnit
	}))
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
