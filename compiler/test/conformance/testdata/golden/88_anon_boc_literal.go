package main

import std "yz/runtime/rt"

type _anonBoc0 struct {
	std.Cown
}


func (self *_anonBoc0) String() string {
	return "_anonBoc0()"
}

func (self *_anonBoc0) describe() std.String {
	return std.NewString("a boc")
}

func (self *_anonBoc0) Describe() *std.Thunk[std.String] {
	return std.Schedule(&self.Cown, func() std.String {
		return self.describe()
	})
}

type Describable interface {
	Describe() *std.Thunk[std.String]
}


type Box[V Describable] struct {
	std.Cown
	value V
}

func NewBox[V Describable](value V) *Box[V] {
	return &Box[V]{
		value: value,
	}
}

func (self *Box[V]) String() string {
	return "Box(" + std.YzTypeName(self.value) + ", " + "value: " + std.StringifyRepr(self.value) + ")"
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
	c := NewBox(&_anonBoc0{})
	std.Print(c.value.Describe().Force())
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
