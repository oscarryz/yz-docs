package main

import std "yz/runtime/rt"

type Wrapper struct {
	std.Cown
	inner *Inner
}

func NewWrapper(inner *Inner) *Wrapper {
	return &Wrapper{
		inner: inner,
	}
}

func (self *Wrapper) String() string {
	return "Wrapper(inner: " + std.StringifyRepr(self.inner) + ")"
}

func (self *Wrapper) Inner() *Inner {
	return self.inner
}

type Inner struct {
	std.Cown
	value std.Int
}

func NewInner(value std.Int) *Inner {
	return &Inner{
		value: value,
	}
}

func (self *Inner) String() string {
	return "Inner(value: " + std.StringifyRepr(self.value) + ")"
}

func (self *Inner) Value() std.Int {
	return self.value
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
	var i *Inner = NewInner(std.NewInt(42))
	var w *Wrapper = NewWrapper(i)
	std.Print(std.NewString(std.StringifyRepr(w.inner.value)))
	return std.TheUnit
}

func (self *_mainBoc) Call() std.Unit {
	return std.LazyUnit(std.Schedule(&self.Cown, func() std.Unit {
		return self.call()
	}))
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
