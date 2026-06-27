package main

import std "yz/runtime/rt"

type _statsBocSwapResult struct {
	_r0 std.String
	_r1 std.String
}


type _statsBoc struct {
	std.Cown
}

func (self *_statsBoc) String() string {
	return "{ " + "swap: {}" + " }"
}

func (self *_statsBoc) swap(a std.String, b std.String) _statsBocSwapResult {
	return _statsBocSwapResult{_r0: b, _r1: a}
}

func (self *_statsBoc) Swap(a std.String, b std.String) *std.Thunk[_statsBocSwapResult] {
	return std.Schedule(&self.Cown, func() _statsBocSwapResult {
		return self.swap(a, b)
	})
}

var Stats = &_statsBoc{}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
	_mrt_x := Stats.Swap(std.NewString("hello"), std.NewString("world")).Force()
	x := _mrt_x._r0
	y := _mrt_x._r1
	std.Print(x)
	std.Print(y)
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
