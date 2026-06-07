package main

import std "yz/runtime/rt"

type _counterBoc struct {
	std.Cown
	count std.Int
}

func (self *_counterBoc) String() string {
	return "{ " + "count: " + std.StringifyRepr(self.count) + "; " + "increment: {}" + "; " + "value: {}" + " }"
}

func (self *_counterBoc) increment(amount std.Int) std.Unit {
	self.count = self.count.Plus(amount)
	return std.TheUnit
}

func (self *_counterBoc) Increment(amount std.Int) std.Unit {
	return std.LazyUnit(std.Schedule(&self.Cown, func() std.Unit {
		return self.increment(amount)
	}))
}

func (self *_counterBoc) value() std.Int {
	return self.count
}

func (self *_counterBoc) Value() std.Int {
	return std.LazyInt(std.Schedule(&self.Cown, func() std.Int {
		return self.value()
	}))
}

var Counter = &_counterBoc{
	count: std.NewInt(0),
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
		var n std.Int
		std.Schedule(&self.Cown, func() std.Unit {
			_st0 := Counter.Increment(std.NewInt(1))
			_bg0.Add(func() { _st0.Await() })
			n = Counter.Value()
			_bg0.Add(func() { n.Await() })
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		_bg1 := &std.BocGroup{}
		_th0 := Counter.Increment(n)
		_bg1.Add(func() { _th0.Await() })
		var m std.Int
		m = Counter.Value()
		_bg1.Add(func() { m.Await() })
		_bg1.Wait()
		std.Print(m.ToStr())
		return std.TheUnit
	}))
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
