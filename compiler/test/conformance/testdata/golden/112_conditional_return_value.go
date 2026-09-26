package main

import std "yz/runtime/rt"

type _pickBoc struct {
	std.Cown
	x std.Int
}

func (self *_pickBoc) String() string {
	return "{ " + "x: " + std.StringifyRepr(self.x) + "; " + "call: {}" + " }"
}

func (self *_pickBoc) Call(x std.Int) std.String {
	return std.LazyString(std.Schedule(&self.Cown, func() std.String {
		self.x = x
		return func() std.String {
			if self.x.Gt(std.NewInt(0)).GoBool() {
				return func() std.String {
					if self.x.Gt(std.NewInt(10)).GoBool() {
						return std.NewString("big")
					} else {
						return std.NewString("small")
					}
				}()
			} else {
				return std.NewString("non-positive")
			}
		}()
	}))
}

var Pick = &_pickBoc{
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
	std.Print(Pick.Call(std.NewInt(20)))
	std.Print(Pick.Call(std.NewInt(5)))
	std.Print(Pick.Call(std.NewInt(1).Neg()))
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
