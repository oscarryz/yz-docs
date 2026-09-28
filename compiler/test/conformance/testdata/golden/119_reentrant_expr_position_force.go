package main

import std "yz/runtime/rt"

type Counter struct {
	std.Cown
}

func NewCounter() *Counter {
	return &Counter{
	}
}

func (self *Counter) String() string {
	return "Counter()"
}

func (self *Counter) countdown(x std.Int) std.Int {
	return func() std.Int {
		if x.Lteq(std.NewInt(0)).GoBool() {
			return std.NewInt(0)
		} else {
			return self.Countdown(x.Minus(std.NewInt(1)))
		}
	}()
}

func (self *Counter) Countdown(x std.Int) std.Int {
	return std.LazyInt(std.Schedule(&self.Cown, func() std.Int {
		return self.countdown(x)
	}))
}

func (self *Counter) Check() std.Bool {
	return std.LazyBool(std.NewThunk(func() std.Bool {
		_bg0 := &std.BocGroup{}
		var __yz_h1 std.Int
		std.Schedule(&self.Cown, func() std.Unit {
			__yz_h1 = self.Countdown(std.NewInt(3))
			_bg0.Add(func() { __yz_h1.Await() })
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		return __yz_h1.Eqeq(std.NewInt(0))
	}))
}

func (self *Counter) Check_arg() std.Bool {
	return std.LazyBool(std.NewThunk(func() std.Bool {
		_bg0 := &std.BocGroup{}
		var __yz_h2 std.Int
		std.Schedule(&self.Cown, func() std.Unit {
			__yz_h2 = self.Countdown(std.NewInt(3))
			_bg0.Add(func() { __yz_h2.Await() })
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		return self.is_zero(__yz_h2)
	}))
}

func (self *Counter) is_zero(n std.Int) std.Bool {
	return n.Eqeq(std.NewInt(0))
}

func (self *Counter) Is_zero(n std.Int) std.Bool {
	return std.LazyBool(std.Schedule(&self.Cown, func() std.Bool {
		return self.is_zero(n)
	}))
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
	var c *Counter = &Counter{}
	std.Print(c.Check().ToStr())
	std.Print(c.Check_arg().ToStr())
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
