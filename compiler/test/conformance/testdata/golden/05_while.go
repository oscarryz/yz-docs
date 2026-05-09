package main

import std "yz/runtime/yzrt"

func while(cond func() std.Bool, body func() std.Unit) *std.Thunk[std.Unit] {
	return std.Go(func() std.Unit {
		if std.Go(func() std.Bool {
			return cond()
		}).GoBool() {
			std.Go(func() std.Unit {
				return body()
			})
			while(cond, body).Force()
		}
		return std.TheUnit
	})
}

type _mainBoc struct {
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.Go(func() std.Unit {
		var n std.Int = std.NewInt(0)
		_bg0 := &std.BocGroup{}
		_bg0.Go(func() any {
			return while(func() std.Bool {
				return n.Lt(std.NewInt(3))
			}, func() std.Unit {
				n = n.Plus(std.NewInt(1))
				return std.TheUnit
			}).Force()
		})
		_bg0.Wait()
		std.Print(n)
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
