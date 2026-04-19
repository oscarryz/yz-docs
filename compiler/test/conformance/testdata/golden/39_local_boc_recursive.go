package main

import std "yz/runtime/yzrt"

type _main_fBoc struct {
}

func (self *_main_fBoc) Call(n std.Int) *std.Thunk[std.Unit] {
	return std.Go(func() std.Unit {
		if n.Eqeq(std.NewInt(0)).GoBool() {
			std.Print(std.NewString("fin"))
		} else {
			std.Print(n)
			self.Call(n.Minus(std.NewInt(1))).Force()
		}
		return std.TheUnit
	})
}


func main() {
	_f := &_main_fBoc{}
	_bg0 := &std.BocGroup{}
	_bg0.Go(func() any {
		return _f.Call(std.NewInt(3)).Force()
	})
	_bg0.Wait()
}
