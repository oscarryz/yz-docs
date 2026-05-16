package main

import std "yz/runtime/rt"

type Box struct {
	std.Cown
	val std.Int
}

func NewBox(val std.Int) *Box {
	return &Box{
		val: val,
	}
}

func (self *Box) set(v std.Int) std.Unit {
	self.val = v
	return std.TheUnit
}

func (self *Box) Set(v std.Int) std.Unit {
	return std.LazyUnit(std.Schedule(&self.Cown, func() std.Unit {
		return self.set(v)
	}))
}

type _applyBoc struct {
	std.Cown
	a *Box
	fn func() std.Unit
}

func (self *_applyBoc) Call(a *Box, fn func() std.Unit) std.Unit {
	return std.LazyUnit(std.ScheduleMulti([]*std.Cown{&self.Cown, &a.Cown}, func() std.Unit {
		self.a = a
		self.fn = fn
		return self.fn()
	}))
}

var Apply = &_applyBoc{
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) Call() std.Unit {
	return std.LazyUnit(std.NewThunk(func() std.Unit {
		_bg0 := &std.BocGroup{}
		var box *Box
		std.Schedule(&self.Cown, func() std.Unit {
			box = NewBox(std.NewInt(0))
			_bg0.GoWait((&_applyBoc{}).Call(box, func() std.Unit {
				return box.set(std.NewInt(42))
			}))
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		std.Print(box.val)
		return std.TheUnit
	}))
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
