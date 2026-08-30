package main

import std "yz/runtime/rt"

type Bag struct {
	std.Cown
	label std.String
	names std.Array[std.String]
}

func NewBag(label std.String, names std.Array[std.String]) *Bag {
	return &Bag{
		label: label,
		names: names,
	}
}

func (self *Bag) String() string {
	return "Bag(label: " + std.StringifyRepr(self.label) + ", names: " + std.StringifyRepr(self.names) + ")"
}

func (self *Bag) Label() std.String {
	return self.label
}

func (self *Bag) Names() std.Array[std.String] {
	return self.names
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
	var b *Bag = NewBag(std.NewString("letters"), std.NewArray(std.NewString("a"), std.NewString("b"), std.NewString("c")))
	std.Print(b.label)
	std.Print(b.names.At(std.NewInt(1)))
	return std.TheUnit
}

func (self *_mainBoc) Call() std.Unit {
	return std.LazyUnit(std.Schedule(&self.Cown, func() std.Unit {
		return self.call()
	}))
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
