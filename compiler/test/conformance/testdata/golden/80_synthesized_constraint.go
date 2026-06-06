package main

import std "yz/runtime/rt"

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

func (self *HolerImpl) Hola() *std.Thunk[std.Unit] {
	return std.Schedule(&self.Cown, func() std.Unit {
		return self.hola()
	})
}

type _WrapperVConstraint interface {
	Hola() *std.Thunk[std.Unit]
}


type Wrapper[V _WrapperVConstraint] struct {
	std.Cown
	item V
}

func NewWrapper[V _WrapperVConstraint](item V) *Wrapper[V] {
	return &Wrapper[V]{
		item: item,
	}
}

func (self *Wrapper[V]) String() string {
	return "Wrapper(" + std.YzTypeName(self.item) + ", " + "item: " + std.StringifyRepr(self.item) + ")"
}

func (self *Wrapper[V]) doIt(value V) std.Unit {
	return value.Hola().Force()
}

func (self *Wrapper[V]) DoIt(value V) *std.Thunk[std.Unit] {
	return std.Schedule(&self.Cown, func() std.Unit {
		return self.doIt(value)
	})
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
		var h *HolerImpl
		var w *Wrapper[*HolerImpl]
		std.Schedule(&self.Cown, func() std.Unit {
			h = &HolerImpl{}
			w = NewWrapper(h)
			_st0 := w.DoIt(h)
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
