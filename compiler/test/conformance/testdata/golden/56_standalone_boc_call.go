package main

import std "yz/runtime/rt"

type _pBoc struct {
	std.Cown
}

func (self *_pBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_pBoc) call() std.Unit {
	std.Print(std.NewString("hello"))
	return std.TheUnit
}

func (self *_pBoc) Call() *std.Thunk[std.Unit] {
	return std.Schedule(&self.Cown, func() std.Unit {
		return self.call()
	})
}

var P = &_pBoc{}

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
			_st0 := P.Call()
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
