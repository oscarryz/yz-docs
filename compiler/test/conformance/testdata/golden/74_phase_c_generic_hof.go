package main

import std "yz/runtime/rt"

func transform[A any, B any](val A, fn func(A) B) *std.Thunk[B] {
	return std.Go(func() B {
		return fn(val)
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

func unwrap[A any](box *Box[A]) *std.Thunk[A] {
	return std.Go(func() A {
		return box.value
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
		std.Schedule(&self.Cown, func() std.Unit {
			std.GoStore(_bg0, transform(std.NewInt(42), func(x std.Int) std.Int {
				return x.Star(std.NewInt(2))
			}), &n)
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		var b *Box[std.Int] = NewBox(std.NewInt(99))
		_bg1 := &std.BocGroup{}
		var v std.Int
		std.GoStore(_bg1, unwrap(b), &v)
		_bg1.Wait()
		std.Print(std.NewString(std.StringifyRepr(n)))
		std.Print(std.NewString(std.StringifyRepr(v)))
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
