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

func (self *Pair[K, V]) String() string {
	return "Pair(" + std.YzTypeName(self.first) + ", " + std.YzTypeName(self.second) + ", " + "first: " + std.StringifyRepr(self.first) + ", second: " + std.StringifyRepr(self.second) + ")"
}

func (self *Pair[K, V]) First() K {
	return self.first
}

func (self *Pair[K, V]) Second() V {
	return self.second
}

func makePair[K any, V any](a K, b V) *std.Thunk[*Pair[K, V]] {
	return std.Go(func() *Pair[K, V] {
		return NewPair(a, b)
	})
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) Call() std.Unit {
	return std.LazyUnit(std.NewThunk(func() std.Unit {
		_bg0 := &std.BocGroup{}
		var p *Pair[std.Int, std.String]
		std.Schedule(&self.Cown, func() std.Unit {
			_st0 := makePair(std.NewInt(42), std.NewString("hello"))
			_bg0.Add(func() { p = _st0.Force() })
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		std.Print(std.NewString(std.StringifyRepr(p.first)))
		std.Print(p.second.ToStr())
		return std.TheUnit
	}))
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
