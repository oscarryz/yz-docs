package main

import std "yz/runtime/rt"

type _count_downBoc struct {
	std.Cown
	n std.Int
}

func (self *_count_downBoc) String() string {
	return "{ " + "n: " + std.StringifyRepr(self.n) + "; " + "call: {}" + " }"
}

func (self *_count_downBoc) Call(n std.Int) std.String {
	return std.LazyString(std.Schedule(&self.Cown, func() std.String {
		self.n = n
		return func() std.String {
			if self.n.Lteq(std.NewInt(0)).GoBool() {
				return std.NewString("done")
			} else {
				return self.Call(self.n.Minus(std.NewInt(1)))
			}
		}()
	}))
}

var Count_down = &_count_downBoc{
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) call() std.Unit {
	std.Print(Count_down.Call(std.NewInt(3)))
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
