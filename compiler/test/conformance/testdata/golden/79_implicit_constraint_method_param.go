package main

import std "yz/runtime/rt"

type Holer interface {
	Hola() std.Unit
}


type HolerImpl struct {
	std.Cown
}

func NewHolerImpl() *HolerImpl {
	return &HolerImpl{
	}
}

func (self *HolerImpl) String() string {
	return "HolerImpl()"
}

func (self *HolerImpl) hola() std.Unit {
	return std.Print(std.NewString("hola"))
}

func (self *HolerImpl) Hola() std.Unit {
	return std.LazyUnit(std.Schedule(&self.Cown, func() std.Unit {
		return self.hola()
	}))
}

type Wrapper[V Holer] struct {
	std.Cown
	item V
}

func NewWrapper[V Holer](item V) *Wrapper[V] {
	return &Wrapper[V]{
		item: item,
	}
}

func (self *Wrapper[V]) String() string {
	return "Wrapper(" + std.YzTypeName(self.item) + ", " + "item: " + std.StringifyRepr(self.item) + ")"
}

func (self *Wrapper[V]) doIt(value V) std.Unit {
	value.Hola().Await()
	return std.TheUnit
}

func (self *Wrapper[V]) DoIt(value V) std.Unit {
	return std.LazyUnit(std.Schedule(&self.Cown, func() std.Unit {
		return self.doIt(value)
	}))
}

func (self *Wrapper[V]) Item() V {
	return self.item
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) Call() std.Unit {
	return std.LazyUnit(std.NewThunk(func() std.Unit {
		_bg0 := &std.BocGroup{}
		var h *HolerImpl
		var w *Wrapper[*HolerImpl]
		std.Schedule(&self.Cown, func() std.Unit {
			h = &HolerImpl{}
			w = NewWrapper(h)
			_st0 := w.DoIt(h)
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
