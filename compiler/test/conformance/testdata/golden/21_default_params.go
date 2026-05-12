package main

import std "yz/runtime/rt"

type _greetBoc struct {
	std.Cown
	name std.String
}

func (self *_greetBoc) Call(name std.String) *std.Thunk[std.Unit] {
	return std.Go(func() std.Unit {
		self.name = name
		return std.Print(self.name)
	})
}

var Greet = &_greetBoc{
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.NewThunk(func() std.Unit {
		_bg0 := &std.BocGroup{}
		std.Schedule(&self.Cown, func() std.Unit {
			_st0 := (&_greetBoc{}).Call(std.NewString("Alice"))
			_bg0.Go(func() any {
				return _st0.Force()
			})
			_st1 := (&_greetBoc{}).Call(std.NewString("Bob"))
			_bg0.Go(func() any {
				return _st1.Force()
			})
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
