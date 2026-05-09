package main

import std "yz/runtime/rt"

type Pair[K any, V any] struct {
	first K
	second V
}

func NewPair[K any, V any](first K, second V) *Pair[K, V] {
	return &Pair[K, V]{
		first: first,
		second: second,
	}
}

func makePair[K any, V any](a K, b V) *std.Thunk[*Pair[K, V]] {
	return std.Go(func() *Pair[K, V] {
		return NewPair(a, b)
	})
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.Schedule(&self.Cown, func() std.Unit {
		p := makePair(std.NewInt(42), std.NewString("hello"))
		std.Print(p.Force().first)
		std.Print(p.Force().second)
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
