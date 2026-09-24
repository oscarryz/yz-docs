package main

import std "yz/runtime/rt"

type _counterBoc struct {
	std.Cown
	count std.Int
}

func (self *_counterBoc) String() string {
	return "{ " + "count: " + std.StringifyRepr(self.count) + "; " + "value: {}" + "; " + "increment: {}" + " }"
}

func (self *_counterBoc) value() std.Int {
	return self.count
}

func (self *_counterBoc) Value() std.Int {
	return std.LazyInt(std.Schedule(&self.Cown, func() std.Int {
		return self.value()
	}))
}

func (self *_counterBoc) increment() std.Unit {
	std.Print(std.NewString("now: ").Plus(self.value().ToStr()))
	self.count = self.count.Plus(std.NewInt(1))
	return std.TheUnit
}

func (self *_counterBoc) Increment() std.Unit {
	return std.LazyUnit(std.Schedule(&self.Cown, func() std.Unit {
		return self.increment()
	}))
}

var Counter = &_counterBoc{
	count: std.NewInt(0),
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "foo: {}" + "; " + "bar: {}" + "; " + "call: {}" + " }"
}

func (self *_mainBoc) foo() std.String {
	return std.NewString("hello")
}

func (self *_mainBoc) Foo() std.String {
	return std.LazyString(std.Schedule(&self.Cown, func() std.String {
		return self.foo()
	}))
}

func (self *_mainBoc) bar() std.String {
	return std.NewString("world")
}

func (self *_mainBoc) Bar() std.String {
	return std.LazyString(std.Schedule(&self.Cown, func() std.String {
		return self.bar()
	}))
}

func (self *_mainBoc) Call() std.Unit {
	return std.LazyUnit(std.NewThunk(func() std.Unit {
		_bg0 := &std.BocGroup{}
		std.Schedule(&self.Cown, func() std.Unit {
			std.Print(self.foo())
			std.Print(self.bar())
			_st0 := Counter.Increment()
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
