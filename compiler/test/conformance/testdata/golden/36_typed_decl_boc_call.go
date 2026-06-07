package main

import std "yz/runtime/rt"

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) Call() std.Unit {
	return std.LazyUnit(std.NewThunk(func() std.Unit {
		_bg0 := &std.BocGroup{}
		var a std.String
		var b std.String
		std.Schedule(&self.Cown, func() std.Unit {
			a = std.Http.Get(std.NewString("https://httpbin.org/get"))
			_bg0.Add(func() { a.Await() })
			b = std.Http.Get(std.NewString("https://httpbin.org/uuid"))
			_bg0.Add(func() { b.Await() })
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		std.Print(a)
		std.Print(b)
		return std.TheUnit
	}))
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
