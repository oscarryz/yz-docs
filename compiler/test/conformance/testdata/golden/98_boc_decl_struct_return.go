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

func (self *Animal) speak() std.String {
	return self.name
}

func (self *Animal) Speak() std.String {
	return std.LazyString(std.Schedule(&self.Cown, func() std.String {
		return self.speak()
	}))
}

func (self *Animal) Name() std.String {
	return self.name
}

type _fooBoc struct {
	std.Cown
}

func (self *_fooBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_fooBoc) Call() std.String {
	return std.LazyString(std.NewThunk(func() std.String {
		var _bocret std.String
		_bg0 := &std.BocGroup{}
		var a *Animal
		std.Schedule(&self.Cown, func() std.Unit {
			a = NewAnimal(std.NewString("cat"))
			_bocret = a.Speak()
			_bg0.Add(func() { _bocret.Await() })
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		return _bocret
	}))
}

var Foo = &_fooBoc{}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
	std.Print(Foo.Call())
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
