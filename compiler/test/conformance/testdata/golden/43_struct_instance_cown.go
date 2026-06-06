package main

import std "yz/runtime/rt"

type Counter struct {
	std.Cown
	count std.Int
}

func NewCounter(count std.Int) *Counter {
	return &Counter{
		count: count,
	}
}

func (self *Counter) String() string {
	return "Counter(count: " + std.StringifyRepr(self.count) + ")"
}

func (self *Counter) increment() std.Unit {
	self.count = self.count.Plus(std.NewInt(1))
	return std.TheUnit
}

func (self *Counter) Increment() *std.Thunk[std.Unit] {
	return std.Schedule(&self.Cown, func() std.Unit {
		return self.increment()
	})
}

func (self *Counter) value() std.Int {
	return self.count
}

func (self *Counter) Value() *std.Thunk[std.Int] {
	return std.Schedule(&self.Cown, func() std.Int {
		return self.value()
	})
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
		var c *Counter
		std.Schedule(&self.Cown, func() std.Unit {
			c = NewCounter(std.NewInt(0))
			_st0 := c.Increment()
			_bg0.Add(func() { _st0.Force() })
			_st1 := c.Increment()
			_bg0.Add(func() { _st1.Force() })
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		std.Print(std.NewString(std.StringifyRepr(c.Value().Force())))
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
