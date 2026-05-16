package main

import std "yz/runtime/rt"

type _counterBoc struct {
	std.Cown
	count std.Int
}

func (self *_counterBoc) increment(amount std.Int) std.Unit {
	std.Print(std.NewString("incrementing ").Plus(std.NewString(std.Stringify(amount))))
	self.count = self.count.Plus(amount)
	return std.TheUnit
}

func (self *_counterBoc) Increment(amount std.Int) *std.Thunk[std.Unit] {
	return std.Schedule(&self.Cown, func() std.Unit {
		return self.increment(amount)
	})
}

func (self *_counterBoc) value() std.Int {
	return self.count
}

func (self *_counterBoc) Value() *std.Thunk[std.Int] {
	return std.Schedule(&self.Cown, func() std.Int {
		return self.value()
	})
}

var Counter = &_counterBoc{
	count: std.NewInt(0),
}

type _pBoc struct {
	std.Cown
}

func (self *_pBoc) call() std.Unit {
	return std.Print(std.NewString("about to print"))
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

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.NewThunk(func() std.Unit {
		_bg0 := &std.BocGroup{}
		std.Schedule(&self.Cown, func() std.Unit {
			_bg0.GoWait(Counter.Increment(std.NewInt(1)))
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		var n std.Int
		_bgs_n := &std.BocGroup{}
		std.GoStore(_bgs_n, Counter.Value(), &n)
		_bgs_n.Wait()
		_bg1 := &std.BocGroup{}
		_bg1.GoWait(Counter.Increment(n))
		_bg1.Wait()
		var m std.Int
		_bgs_m := &std.BocGroup{}
		std.GoStore(_bgs_m, Counter.Value(), &m)
		_bgs_m.Wait()
		_bg2 := &std.BocGroup{}
		_bg2.GoWait(P.Call())
		_bg2.Wait()
		std.Print(std.NewString(std.Stringify(m)))
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
