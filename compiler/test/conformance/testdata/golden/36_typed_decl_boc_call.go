package main

import std "yz/runtime/rt"

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.NewThunk(func() std.Unit {
		var a std.String
		_bgs_a := &std.BocGroup{}
		std.Schedule(&self.Cown, func() std.Unit {
			std.GoStore(_bgs_a, std.Http.Get(std.NewString("https://httpbin.org/get")), &a)
			return std.TheUnit
		}).Force()
		_bgs_a.Wait()
		var b std.String
		_bgs_b := &std.BocGroup{}
		std.GoStore(_bgs_b, std.Http.Get(std.NewString("https://httpbin.org/uuid")), &b)
		_bgs_b.Wait()
		std.Print(a)
		std.Print(b)
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
