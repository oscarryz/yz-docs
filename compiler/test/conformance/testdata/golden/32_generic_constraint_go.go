package main

import std "yz/runtime/yzrt"

type Wrapper[T interface{ ToStr() std.String }] struct {
	value T
}

func NewWrapper[T interface{ ToStr() std.String }](value T) *Wrapper[T] {
	return &Wrapper[T]{
		value: value,
	}
}

func (self *Wrapper[T]) Describe() *std.Thunk[std.String] {
	return std.Go(func() std.String {
		return self.value.ToStr()
	})
}

func main() {
	w := NewWrapper(std.NewString("hello"))
	std.Print(w.Describe().Force())
}
