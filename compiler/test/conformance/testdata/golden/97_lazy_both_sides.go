package main

import std "yz/runtime/rt"

type _someBoc struct {
	std.Cown
}

func (self *_someBoc) String() string {
	return "{ " + "get: {}" + " }"
}

func (self *_someBoc) get() std.String {
	return std.NewString("hello")
}

func (self *_someBoc) Get() std.String {
	return std.LazyString(std.Schedule(&self.Cown, func() std.String {
		return self.get()
	}))
}

var Some = &_someBoc{}

type _moreBoc struct {
	std.Cown
}

func (self *_moreBoc) String() string {
	return "{ " + "get: {}" + " }"
}

func (self *_moreBoc) get() std.String {
	return std.NewString("hello")
}

func (self *_moreBoc) Get() std.String {
	return std.LazyString(std.Schedule(&self.Cown, func() std.String {
		return self.get()
	}))
}

var More = &_moreBoc{}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
	if Some.Get().Eqeq(More.Get()).GoBool() {
		std.Print(std.NewString("equal"))
	} else {
		std.Print(std.NewString("not equal"))
	}
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
