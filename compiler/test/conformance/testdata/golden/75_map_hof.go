package main

import std "yz/runtime/rt"

func mapList[A any, B any](list std.Array[A], fn func(A) B) *std.Thunk[std.Array[B]] {
	return std.Go(func() std.Array[B] {
		return std.ArrayMap(list, func(item A) B {
			return fn(item)
		})
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
		var doubled std.Array[std.Int]
		std.Schedule(&self.Cown, func() std.Unit {
			std.GoStore(_bg0, mapList(std.NewArray(std.NewInt(1), std.NewInt(2), std.NewInt(3)), func(x std.Int) std.Int {
				return x.Star(std.NewInt(2))
			}), &doubled)
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		std.Print(doubled)
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
