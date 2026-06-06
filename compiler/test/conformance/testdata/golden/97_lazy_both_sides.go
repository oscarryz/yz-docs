package main

import std "yz/runtime/rt"

type _someBoc struct {
	std.Cown
}

func (self *_someBoc) String() string {
	return "{ " + "get: {}" + " }"
}

func (self *_someBoc) get() std.String {
	return std.NewString("hello")
}

func (self *_someBoc) Get() *std.Thunk[std.String] {
	return std.Schedule(&self.Cown, func() std.String {
		return self.get()
	})
}

var Some = &_someBoc{}

type _moreBoc struct {
	std.Cown
}

func (self *_moreBoc) String() string {
	return "{ " + "get: {}" + " }"
}

func (self *_moreBoc) get() std.String {
	return std.NewString("hello")
}

func (self *_moreBoc) Get() *std.Thunk[std.String] {
	return std.Schedule(&self.Cown, func() std.String {
		return self.get()
	})
}

var More = &_moreBoc{}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.NewThunk(func() std.Unit {
		_bg0 := &std.BocGroup{}
		std.Schedule(&self.Cown, func() std.Unit {
			_st0 := std.WrapStringThunk(Some.Get()).EqeqT(std.WrapStringThunk(More.Get())).Qm(func() any {
				std.Print(std.NewString("equal"))
				return std.TheUnit
			}, func() any {
				std.Print(std.NewString("not equal"))
				return std.TheUnit
			})
			_bg0.Add(func() { _st0.Force() })
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
