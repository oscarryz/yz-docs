package main

import std "yz/runtime/rt"

type _countdownBoc struct {
	std.Cown
	n std.Int
}

func (self *_countdownBoc) String() string {
	return "{ " + "n: " + std.StringifyRepr(self.n) + "; " + "call: {}" + " }"
}

func (self *_countdownBoc) Call(n std.Int) std.Unit {
	return std.LazyUnit(std.NewThunk(func() std.Unit {
		_bg0 := &std.BocGroup{}
		std.Schedule(&self.Cown, func() std.Unit {
			self.n = n
			if self.n.Eqeq(std.NewInt(0)).GoBool() {
				std.Print(std.NewString("done"))
			} else {
				std.Print(std.NewString(std.StringifyRepr(self.n)))
				_st0 := self.Call(self.n.Minus(std.NewInt(1)))
				_bg0.Add(func() { _st0.Await() })
			}
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		return std.TheUnit
	}))
}

var Countdown = &_countdownBoc{
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) Call() std.Unit {
	return std.LazyUnit(std.NewThunk(func() std.Unit {
		_bg0 := &std.BocGroup{}
		std.Schedule(&self.Cown, func() std.Unit {
			_st0 := Countdown.Call(std.NewInt(3))
			_bg0.Add(func() { _st0.Await() })
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		return std.TheUnit
	}))
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
