package main

import std "yz/runtime/rt"

type _swapBocResult struct {
	_r0 std.String
	_r1 std.String
}


type _swapBoc struct {
	std.Cown
}

func (self *_swapBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_swapBoc) call(a std.String, b std.String) _swapBocResult {
	return _swapBocResult{_r0: b, _r1: a}
}

func (self *_swapBoc) Call(a std.String, b std.String) *std.Thunk[_swapBocResult] {
	return std.Schedule(&self.Cown, func() _swapBocResult {
		return self.call(a, b)
	})
}

var Swap = &_swapBoc{}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
	_mrt_x := Swap.Call(std.NewString("hello"), std.NewString("world")).Force()
	x := _mrt_x._r0
	y := _mrt_x._r1
	std.Print(x)
	std.Print(y)
	return std.TheUnit
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.Schedule(&self.Cown, func() std.Unit {
		return self.call()
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
