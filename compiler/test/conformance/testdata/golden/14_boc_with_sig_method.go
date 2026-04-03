package main

import std "yz/runtime/yzrt"

type _counterBoc struct {
	count std.Int
}

func (self *_counterBoc) increment() *std.Thunk[std.Unit] {
	return std.Go(func() std.Unit {
		self.count = self.count.Plus(std.NewInt(1))
		return std.TheUnit
	})
}

func (self *_counterBoc) value() *std.Thunk[std.Int] {
	return std.Go(func() std.Int {
		return self.count
	})
}

var counter = &_counterBoc{
	count: std.NewInt(0),
}
