package main

import std "yz/runtime/rt"

type Field struct {
	std.Cown
	name std.String
	type_ std.String
}

func NewField(name std.String, type_ std.String) *Field {
	return &Field{
		name: name,
		type_: type_,
	}
}

func (self *Field) String() string {
	return "Field(name: " + std.StringifyRepr(self.name) + ", type: " + std.StringifyRepr(self.type_) + ")"
}

func (self *Field) Name() std.String {
	return self.name
}

func (self *Field) Type() std.String {
	return self.type_
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
	var f *Field = NewField(std.NewString("age"), std.NewString("Int"))
	std.Print(f.type_)
	std.Print(std.NewString(std.StringifyRepr(f)))
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
