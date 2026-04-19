package main

import std "yz/runtime/yzrt"

func greet(name std.String) *std.Thunk[std.Unit] {
	return std.Go(func() std.Unit {
		return std.Print(name)
	})
}

type _mainBoc struct {
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.Go(func() std.Unit {
		_bg0 := &std.BocGroup{}
		_bg0.Go(func() any {
			return greet(std.NewString("Alice")).Force()
		})
		_bg0.Go(func() any {
			return greet(std.NewString("Bob")).Force()
		})
		_bg0.Wait()
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
