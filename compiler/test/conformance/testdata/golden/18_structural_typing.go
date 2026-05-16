package main

import std "yz/runtime/rt"

type Greeter interface {
	greet() *std.Thunk[std.Unit]
}


type Person struct {
	std.Cown
	name std.String
	secret std.String
}

func NewPerson(name std.String, secret std.String) *Person {
	return &Person{
		name: name,
		secret: secret,
	}
}

func (self *Person) greet() std.Unit {
	return std.Print(self.name)
}

func (self *Person) Greet() *std.Thunk[std.Unit] {
	return std.Schedule(&self.Cown, func() std.Unit {
		return self.greet()
	})
}

type _greet_allBoc struct {
	std.Cown
	g Greeter
}

func (self *_greet_allBoc) Call(g Greeter) *std.Thunk[std.Unit] {
	return std.Go(func() std.Unit {
		self.g = g
		return self.g.Greet().Force()
	})
}

var Greet_all = &_greet_allBoc{
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.NewThunk(func() std.Unit {
		_bg0 := &std.BocGroup{}
		var p *Person
		std.Schedule(&self.Cown, func() std.Unit {
			p = NewPerson(std.NewString("Alice"), std.NewString("my secret"))
			_bg0.GoWait((&_greet_allBoc{}).Call(p))
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
