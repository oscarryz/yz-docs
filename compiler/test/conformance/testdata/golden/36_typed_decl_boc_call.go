package main

import std "yz/runtime/rt"

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) call() std.Unit {
	a := std.Http.Get(std.NewString("https://httpbin.org/get"))
	b := std.Http.Get(std.NewString("https://httpbin.org/uuid"))
	std.Print(a)
	std.Print(b)
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
