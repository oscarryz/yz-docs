package main

import std "yz/runtime/rt"

type _counterBoc struct {
	std.Cown
	count std.Int
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
