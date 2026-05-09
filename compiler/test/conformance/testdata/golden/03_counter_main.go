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

func (self *_counterBoc) Value() *std.Thunk[std.Int] {
	return std.Go(func() std.Int {
		return self.count
	})
}

var Counter = &_counterBoc{
	count: std.NewInt(0),
}

type _mainBoc struct {
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.Go(func() std.Unit {
		_bg0 := &std.BocGroup{}
		_bg0.Go(func() any {
			return Counter.Increment().Force()
		})
		_bg0.Go(func() any {
			return Counter.Increment().Force()
		})
		_bg0.Wait()
		std.Print(Counter.Value().Force())
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
