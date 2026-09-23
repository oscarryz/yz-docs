package main

import std "yz/runtime/rt"

type _my_funcBoc struct {
	std.Cown
}

func (self *_my_funcBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_my_funcBoc) call(name std.String, greet func(std.String) std.Unit) std.Unit {
	std.Print(name)
	greet(std.NewString("Hello"))
	return std.TheUnit
}

func (self *_my_funcBoc) Call(name std.String, greet func(std.String) std.Unit) std.Unit {
	return std.LazyUnit(std.Schedule(&self.Cown, func() std.Unit {
		return self.call(name, greet)
	}))
}

var My_func = &_my_funcBoc{}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) Call() std.Unit {
	return std.LazyUnit(std.NewThunk(func() std.Unit {
		_bg0 := &std.BocGroup{}
		std.Schedule(&self.Cown, func() std.Unit {
			_st0 := My_func.Call(std.NewString("World"), func(msg std.String) std.Unit {
				return std.Print(msg)
			})
			_bg0.Add(func() { _st0.Await() })
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		return std.TheUnit
	}))
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
