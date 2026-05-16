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
	std.Print(b.value)
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
