package main

import std "yz/runtime/rt"

type Pair[K any, V any] struct {
	std.Cown
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
	return std.NewThunk(func() std.Unit {
		var p *Pair[K, V]
		_bgs_p := &std.BocGroup{}
		std.Schedule(&self.Cown, func() std.Unit {
			std.GoStore(_bgs_p, makePair(std.NewInt(42), std.NewString("hello")), &p)
			return std.TheUnit
		}).Force()
		_bgs_p.Wait()
		std.Print(p.first)
		std.Print(p.second)
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
