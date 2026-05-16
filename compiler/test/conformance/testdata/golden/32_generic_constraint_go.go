package main

import std "yz/runtime/rt"

type Wrapper[T interface{ ToStr() std.String }] struct {
	std.Cown
	value T
}

func NewWrapper[T interface{ ToStr() std.String }](value T) *Wrapper[T] {
	return &Wrapper[T]{
		value: value,
	}
}

func (self *Wrapper[T]) describe() std.String {
	return self.value.ToStr()
}

func (self *Wrapper[T]) Describe() std.String {
	return std.LazyString(std.Schedule(&self.Cown, func() std.String {
		return self.describe()
	}))
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) call() std.Unit {
	w := NewWrapper(std.NewString("hello"))
	std.Print(w.Describe())
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
