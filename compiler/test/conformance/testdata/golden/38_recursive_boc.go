package main

import std "yz/runtime/rt"

type _countdownBoc struct {
	std.Cown
	n std.Int
}

func (self *_countdownBoc) Call(n std.Int) *std.Thunk[std.Unit] {
	return std.Go(func() std.Unit {
		self.n = n
		if self.n.Eqeq(std.NewInt(0)).GoBool() {
			std.Print(std.NewString("done"))
		} else {
			std.Print(self.n)
			(&_countdownBoc{}).Call(self.n.Minus(std.NewInt(1))).Force()
		}
		return std.TheUnit
	})
}

var Countdown = &_countdownBoc{
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.NewThunk(func() std.Unit {
		_bg0 := &std.BocGroup{}
		std.Schedule(&self.Cown, func() std.Unit {
			_st0 := (&_countdownBoc{}).Call(std.NewInt(3))
			_bg0.Go(func() any {
				return _st0.Force()
			})
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
