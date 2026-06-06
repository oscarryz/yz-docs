package main

import std "yz/runtime/rt"

type _counterBoc struct {
	std.Cown
	n std.Int
}

func (self *_counterBoc) String() string {
	return "{ " + "n: " + std.StringifyRepr(self.n) + "; " + "increment: {}" + "; " + "value: {}" + " }"
}

func (self *_counterBoc) increment() std.Unit {
	self.n = self.n.Plus(std.NewInt(1))
	return std.TheUnit
}

func (self *_counterBoc) Increment() *std.Thunk[std.Unit] {
	return std.Schedule(&self.Cown, func() std.Unit {
		return self.increment()
	})
}

func (self *_counterBoc) value() std.Int {
	return self.n
}

func (self *_counterBoc) Value() *std.Thunk[std.Int] {
	return std.Schedule(&self.Cown, func() std.Int {
		return self.value()
	})
}

var Counter = &_counterBoc{
	n: std.NewInt(0),
}

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
			_st0 := Counter.Increment()
			_bg0.Add(func() { _st0.Force() })
			_st1 := Counter.Increment()
			_bg0.Add(func() { _st1.Force() })
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		std.Print(std.NewString(std.StringifyRepr(Counter.Value().Force())))
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
