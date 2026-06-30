package main

import std "yz/runtime/rt"

type Animal struct {
	std.Cown
	name std.String
}

func NewAnimal(name std.String) *Animal {
	return &Animal{
		name: name,
	}
}

func (self *Animal) String() string {
	return "Animal(name: " + std.StringifyRepr(self.name) + ")"
}

func (self *Animal) describe() std.String {
	return self.name
}

func (self *Animal) Describe() std.String {
	return std.LazyString(std.Schedule(&self.Cown, func() std.String {
		return self.describe()
	}))
}

func (self *Animal) Name() std.String {
	return self.name
}

type _BoxVConstraint interface {
	Describe() std.String
}


type Box[V _BoxVConstraint] struct {
	std.Cown
	value V
}

func NewBox[V _BoxVConstraint](value V) *Box[V] {
	return &Box[V]{
		value: value,
	}
}

func (self *Box[V]) String() string {
	return "Box(" + std.YzTypeName(self.value) + ", " + "value: " + std.StringifyRepr(self.value) + ")"
}

func (self *Box[V]) desc() std.String {
	return self.value.Describe()
}

func (self *Box[V]) Desc() std.String {
	return std.LazyString(std.Schedule(&self.Cown, func() std.String {
		return self.desc()
	}))
}

func (self *Box[V]) Value() V {
	return self.value
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
	var a *Animal = NewAnimal(std.NewString("Cat"))
	var b *Box[*Animal] = NewBox(a)
	std.Print(b.Desc())
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
