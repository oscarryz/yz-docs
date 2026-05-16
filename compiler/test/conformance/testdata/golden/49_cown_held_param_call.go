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

func (self *Box) Set(v std.Int) *std.Thunk[std.Unit] {
	return std.Schedule(&self.Cown, func() std.Unit {
		return self.set(v)
	})
}

type _assignBoc struct {
	std.Cown
	b *Box
	v std.Int
}

func (self *_assignBoc) Call(b *Box, v std.Int) *std.Thunk[std.Unit] {
	return func() *std.Thunk[std.Unit] {
		_bg0 := &std.BocGroup{}
		_sched := std.ScheduleMulti([]*std.Cown{&self.Cown, &b.Cown}, func() std.Unit {
			self.b = b
			self.v = v
			_bg0.GoWait(self.b.Set(self.v))
			return std.TheUnit
		})
		return std.NewThunk(func() std.Unit {
			_sched.Force()
			_bg0.Wait()
			return std.TheUnit
		})
	}()
}

var Assign = &_assignBoc{
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.NewThunk(func() std.Unit {
		_bg0 := &std.BocGroup{}
		var box *Box
		std.Schedule(&self.Cown, func() std.Unit {
			box = NewBox(std.NewInt(0))
			_bg0.GoWait((&_assignBoc{}).Call(box, std.NewInt(42)))
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		std.Print(box.val)
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
