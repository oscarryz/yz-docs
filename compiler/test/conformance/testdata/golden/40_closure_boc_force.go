package main

import std "yz/runtime/rt"

type _counterBoc struct {
	count std.Int
}

func (self *_counterBoc) Increment() *std.Thunk[std.Unit] {
	return std.Go(func() std.Unit {
		self.count = self.count.Plus(std.NewInt(1))
		return std.TheUnit
	})
}

var Counter = &_counterBoc{
	count: std.NewInt(0),
}

type _mainBoc struct {
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.Go(func() std.Unit {
		var list std.Array[std.Int] = std.NewArray(std.NewInt(1), std.NewInt(2), std.NewInt(3))
		list.Each(func(item std.Int) std.Unit {
			Counter.Increment().Force()
			return std.Print(item)
		})
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
