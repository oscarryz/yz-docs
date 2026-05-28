package main

import std "yz/runtime/rt"

type Container[T any] struct {
	std.Cown
	value T
}

func NewContainer[T any](value T) *Container[T] {
	return &Container[T]{
		value: value,
	}
}

func (self *Container[T]) String() string {
	return "Container(" + std.YzTypeName(self.value) + ", " + "value: " + std.StringifyRepr(self.value) + ")"
}

func (self *Container[T]) get() T {
	return self.value
}

func (self *Container[T]) Get() *std.Thunk[T] {
	return std.Schedule(&self.Cown, func() T {
		return self.get()
	})
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
	var c *Container[std.Int] = NewContainer(std.NewInt(42))
	var s *Container[std.String] = NewContainer(std.NewString("hello"))
	std.Print(c.Get().Force())
	std.Print(s.Get().Force())
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
