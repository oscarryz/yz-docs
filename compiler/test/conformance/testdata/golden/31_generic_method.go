package main

import std "yz/runtime/yzrt"

type Container[T any] struct {
	value T
}

func NewContainer[T any](value T) *Container[T] {
	return &Container[T]{
		value: value,
	}
}

func (self *Container[T]) Get() *std.Thunk[T] {
	return std.Go(func() T {
		return self.value
	})
}

func main() {
	c := NewContainer(std.NewInt(42))
	s := NewContainer(std.NewString("hello"))
	std.Print(c.Get().Force())
	std.Print(s.Get().Force())
}
