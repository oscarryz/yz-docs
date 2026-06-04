package main

import std "yz/runtime/rt"

func identity[A any](val A) *std.Thunk[A] {
	return std.Go(func() A {
		return val
	})
}

type Box[T any] struct {
	std.Cown
	value T
}

func NewBox[T any](value T) *Box[T] {
	return &Box[T]{
		value: value,
	}
}

func (self *Box[T]) String() string {
	return "Box(" + std.YzTypeName(self.value) + ", " + "value: " + std.StringifyRepr(self.value) + ")"
}

func wrap[A any](val A) *std.Thunk[*Box[A]] {
	return std.Go(func() *Box[A] {
		return NewBox(val)
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
		var n std.Int
		var s std.String
		var b *Box[std.Int]
		std.Schedule(&self.Cown, func() std.Unit {
			std.GoStore(_bg0, identity(std.NewInt(42)), &n)
			std.GoStore(_bg0, identity(std.NewString("hello")), &s)
			std.GoStore(_bg0, wrap(std.NewInt(99)), &b)
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		std.Print(std.NewString(std.StringifyRepr(n)))
		std.Print(s)
		std.Print(std.NewString(std.StringifyRepr(b.value)))
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
