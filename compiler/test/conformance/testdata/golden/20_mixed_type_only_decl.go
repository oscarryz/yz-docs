package main

import std "yz/runtime/rt"

type Named struct {
	name std.String
	greet func() *std.Thunk[std.Unit]
}

func NewNamed(name std.String, greet func() *std.Thunk[std.Unit]) *Named {
	return &Named{
		name: name,
		greet: greet,
	}
}

func (self *Named) Greet() *std.Thunk[std.Unit] {
	return self.greet()
}

type _mainBoc struct {
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.Go(func() std.Unit {
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
