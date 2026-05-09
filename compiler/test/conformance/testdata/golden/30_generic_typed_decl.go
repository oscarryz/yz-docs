package main

import std "yz/runtime/rt"

type Box[T any] struct {
	value T
}

func NewBox[T any](value T) *Box[T] {
	return &Box[T]{
		value: value,
	}
}

type _mainBoc struct {
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.Go(func() std.Unit {
		b := NewBox(std.NewInt(42))
		var s *Box[std.String] = NewBox(std.NewString("hello"))
		std.Print(b.value)
		std.Print(s.value)
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
