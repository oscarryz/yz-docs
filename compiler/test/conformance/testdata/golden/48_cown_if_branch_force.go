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

type _cond_setBoc struct {
	std.Cown
	a *Box
	flag std.Bool
}

func (self *_cond_setBoc) Call(a *Box, flag std.Bool) std.Unit {
	return std.LazyUnit(func() *std.Thunk[std.Unit] {
		_bg0 := &std.BocGroup{}
		_sched := std.ScheduleMulti([]*std.Cown{&self.Cown, &a.Cown}, func() std.Unit {
			self.a = a
			self.flag = flag
			if self.flag.GoBool() {
				_bg0.GoWait(self.a.Set(std.NewInt(1)))
			} else {
				_bg0.GoWait(self.a.Set(std.NewInt(0)))
			}
			return std.TheUnit
		})
		return std.NewThunk(func() std.Unit {
			_sched.Force()
			_bg0.Wait()
			return std.TheUnit
		})
	}())
}

var Cond_set = &_cond_setBoc{
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) Call() std.Unit {
	return std.LazyUnit(std.NewThunk(func() std.Unit {
		_bg0 := &std.BocGroup{}
		var b *Box
		std.Schedule(&self.Cown, func() std.Unit {
			b = NewBox(std.NewInt(0))
			_bg0.GoWait((&_cond_setBoc{}).Call(b, std.NewBool(true)))
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		std.Print(b.val)
		return std.TheUnit
	}))
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
