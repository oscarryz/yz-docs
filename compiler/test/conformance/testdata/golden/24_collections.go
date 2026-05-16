package main

import std "yz/runtime/rt"

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) call() std.Unit {
	var nums std.Array[std.Int] = std.NewArray(std.NewInt(1), std.NewInt(2), std.NewInt(3))
	var scores std.Dict[std.String, std.Int] = std.NewDict[std.String, std.Int]().Set(std.NewString("alice"), std.NewInt(10)).Set(std.NewString("bob"), std.NewInt(20))
	std.Print(nums.At(std.NewInt(0)))
	std.Print(scores.At(std.NewString("alice")))
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
