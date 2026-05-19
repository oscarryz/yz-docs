package main

import std "yz/runtime/rt"

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.NewThunk(func() std.Unit {
		_bg0 := &std.BocGroup{}
		var a std.String
		var b std.String
		std.Schedule(&self.Cown, func() std.Unit {
			std.GoStore(_bg0, std.Http.Get(std.NewString("https://httpbin.org/get")), &a)
			std.GoStore(_bg0, std.Http.Get(std.NewString("https://httpbin.org/uuid")), &b)
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		std.Print(a)
		std.Print(b)
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
