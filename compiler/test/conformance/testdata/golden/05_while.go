package main

import std "yz/runtime/rt"

type _whileBoc struct {
	std.Cown
	cond func() std.Bool
	body func() std.Unit
}

func (self *_whileBoc) String() string {
	return "{ " + "cond: " + std.StringifyRepr(self.cond) + "; " + "body: " + std.StringifyRepr(self.body) + "; " + "call: {}" + " }"
}

func (self *_whileBoc) Call(cond func() std.Bool, body func() std.Unit) *std.Thunk[std.Unit] {
	return std.NewThunk(func() std.Unit {
		_bg0 := &std.BocGroup{}
		std.Schedule(&self.Cown, func() std.Unit {
			self.cond = cond
			self.body = body
			if self.cond().GoBool() {
				self.body()
				_bg0.GoWait(self.Call(self.cond, self.body))
			}
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		return std.TheUnit
	})
}

var While = &_whileBoc{
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
		var n std.Int
		std.Schedule(&self.Cown, func() std.Unit {
			n = std.NewInt(0)
			_bg0.GoWait(While.Call(func() std.Bool {
				return n.Lt(std.NewInt(3))
			}, func() std.Unit {
				n = n.Plus(std.NewInt(1))
				return std.TheUnit
			}))
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		std.Print(std.NewString(std.StringifyRepr(n)))
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
