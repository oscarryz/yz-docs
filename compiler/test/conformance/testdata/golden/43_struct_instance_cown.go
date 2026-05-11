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

func (self *Counter) Increment() *std.Thunk[std.Unit] {
	return std.Schedule(&self.Cown, func() std.Unit {
		self.count = self.count.Plus(std.NewInt(1))
		return std.TheUnit
	})
}

func (self *Counter) Value() *std.Thunk[std.Int] {
	return std.Schedule(&self.Cown, func() std.Int {
		return self.count
	})
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.NewThunk(func() std.Unit {
		_bg0 := &std.BocGroup{}
		var c *Counter
		std.Schedule(&self.Cown, func() std.Unit {
			c = NewCounter(std.NewInt(0))
			_bg0.Go(func() any {
				return c.Increment().Force()
			})
			_bg0.Go(func() any {
				return c.Increment().Force()
			})
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		std.Print(c.Value().Force())
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
