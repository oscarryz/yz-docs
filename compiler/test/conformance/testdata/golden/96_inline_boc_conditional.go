package main

import std "yz/runtime/rt"

type _greeterBoc struct {
	std.Cown
}

func (self *_greeterBoc) String() string {
	return "{ " + "greet: {}" + " }"
}

func (self *_greeterBoc) greet() std.String {
	return std.NewString("hello")
}

func (self *_greeterBoc) Greet() *std.Thunk[std.String] {
	return std.Schedule(&self.Cown, func() std.String {
		return self.greet()
	})
}

var Greeter = &_greeterBoc{}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
	if Greeter.Greet().Force().Eqeq(std.NewString("hello")).GoBool() {
		std.Print(std.NewString("matched"))
	} else {
		std.Print(std.NewString("no match"))
	}
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
