package main

import std "yz/runtime/rt"

type _greetBoc struct {
	std.Cown
}

func (self *_greetBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_greetBoc) call(name std.String) std.Unit {
	std.Print(std.NewString("Hello, ").Plus(name.ToStr()).Plus(std.NewString("!")))
	return std.TheUnit
}

func (self *_greetBoc) Call(name std.String) *std.Thunk[std.Unit] {
	return std.Schedule(&self.Cown, func() std.Unit {
		return self.call(name)
	})
}

var Greet = &_greetBoc{}

type _addBoc struct {
	std.Cown
}

func (self *_addBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_addBoc) call(a std.Int, b std.Int) std.Int {
	return a.Plus(b)
}

func (self *_addBoc) Call(a std.Int, b std.Int) *std.Thunk[std.Int] {
	return std.Schedule(&self.Cown, func() std.Int {
		return self.call(a, b)
	})
}

var Add = &_addBoc{}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.NewThunk(func() std.Unit {
		_bg0 := &std.BocGroup{}
		var result std.Int
		std.Schedule(&self.Cown, func() std.Unit {
			_st0 := Greet.Call(std.NewString("World"))
			_bg0.Add(func() { _st0.Force() })
			_st1 := Add.Call(std.NewInt(3), std.NewInt(4))
			_bg0.Add(func() { result = _st1.Force() })
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		std.Print(std.NewString(std.StringifyRepr(result)))
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
