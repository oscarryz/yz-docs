package main

import std "yz/runtime/rt"

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) call() std.Unit {
	var score std.Int = std.NewInt(85)
	var grade std.String = func() std.String {
		if score.Gteq(std.NewInt(90)).GoBool() {
			return std.NewString("A")
		} else if score.Gteq(std.NewInt(80)).GoBool() {
			return std.NewString("B")
		} else if score.Gteq(std.NewInt(70)).GoBool() {
			return std.NewString("C")
		} else {
			return std.NewString("F")
		}
	}()
	std.Print(grade)
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
