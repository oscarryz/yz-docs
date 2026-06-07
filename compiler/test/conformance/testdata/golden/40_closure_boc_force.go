package main

import std "yz/runtime/rt"

type _counterBoc struct {
	std.Cown
	count std.Int
}

func (self *_counterBoc) String() string {
	return "{ " + "count: " + std.StringifyRepr(self.count) + "; " + "increment: {}" + " }"
}

func (self *_counterBoc) increment() std.Unit {
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
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
	var list std.Array[std.Int] = std.NewArray(std.NewInt(1), std.NewInt(2), std.NewInt(3))
	list.Each(func(item std.Int) std.Unit {
		Counter.Increment().Await()
		return std.Print(std.NewString(std.StringifyRepr(item)))
	})
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
