package main

import std "yz/runtime/rt"

type Foo struct {
	std.Cown
	bar std.String
}

func NewFoo(bar std.String) *Foo {
	return &Foo{
		bar: bar,
	}
}

func (self *Foo) String() string {
	return "Foo(bar: " + std.StringifyRepr(self.bar) + ")"
}

func (self *Foo) Bar() std.String {
	return self.bar
}

type _HasBarSchemaBound interface {
	Bar() std.String
}


type HasBar interface {
}


type Quz struct {
	std.Cown
}

func NewQuz() *Quz {
	return &Quz{
	}
}

func (self *Quz) String() string {
	return "Quz()"
}

type _helloBoc struct {
	std.Cown
	q HasBar
	s _HasBarSchemaBound
}

func (self *_helloBoc) String() string {
	return "{ " + "q: " + std.StringifyRepr(self.q) + "; " + "s: " + std.StringifyRepr(self.s) + "; " + "call: {}" + " }"
}

func (self *_helloBoc) Call(q HasBar, s _HasBarSchemaBound) std.String {
	return std.LazyString(std.Schedule(&self.Cown, func() std.String {
		self.q = q
		self.s = s
		return std.NewString("Hello ").Plus(self.s.Bar().ToStr())
	}))
}

var Hello = &_helloBoc{
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
		var result std.String
		var f *Foo
		var q *Quz
		std.Schedule(&self.Cown, func() std.Unit {
			f = NewFoo(std.NewString("world"))
			q = &Quz{}
			result = Hello.Call(q, f)
			_bg0.Add(func() { result.Await() })
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		std.Print(result)
		return std.TheUnit
	}))
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
