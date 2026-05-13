package main

import std "yz/runtime/rt"

type Box[T any] struct {
	std.Cown
	value T
}

func NewBox[T any](value T) *Box[T] {
	return &Box[T]{
		value: value,
	}
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) call() std.Unit {
	b := NewBox(std.NewInt(42))
	var s *Box[std.String] = NewBox(std.NewString("hello"))
	std.Print(b.value)
	std.Print(s.value)
	return std.TheUnit
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.Schedule(&self.Cown, func() std.Unit {
		return self.call()
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
