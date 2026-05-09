package main

import std "yz/runtime/rt"

type Wrapper[T interface{ ToStr() std.String }] struct {
	value T
}

func NewWrapper[T interface{ ToStr() std.String }](value T) *Wrapper[T] {
	return &Wrapper[T]{
		value: value,
	}
}

func (self *Wrapper[T]) Describe() *std.Thunk[std.String] {
	return std.Schedule(&self.Cown, func() std.String {
		return self.value.ToStr()
	})
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.Schedule(&self.Cown, func() std.Unit {
		w := NewWrapper(std.NewString("hello"))
		std.Print(w.Describe().Force())
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
