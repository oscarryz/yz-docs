package main

import std "yz/runtime/rt"

type Printer struct {
	std.Cown
}

func NewPrinter() *Printer {
	return &Printer{
	}
}

func (self *Printer) String() string {
	return "Printer()"
}

func (self *Printer) print_all(items std.Array[std.Int]) std.Unit {
	return items.Each(func(v std.Int) std.Unit {
		return std.Print(v.ToStr())
	})
}

func (self *Printer) Print_all(items std.Array[std.Int]) std.Unit {
	return std.LazyUnit(std.Schedule(&self.Cown, func() std.Unit {
		return self.print_all(items)
	}))
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) String() string {
	return "{ " + "call: {}" + " }"
}

func (self *_mainBoc) Call() std.Unit {
	return std.LazyUnit(std.NewThunk(func() std.Unit {
		_bg0 := &std.BocGroup{}
		var p *Printer
		std.Schedule(&self.Cown, func() std.Unit {
			p = &Printer{}
			_st0 := p.Print_all(std.NewArray(std.NewInt(10), std.NewInt(20), std.NewInt(30)))
			_bg0.Add(func() { _st0.Await() })
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		var d std.Dict[std.String, std.Int] = std.NewDict[std.String, std.Int]().Set(std.NewString("a"), std.NewInt(1)).Set(std.NewString("b"), std.NewInt(2))
		d.Each(func(k std.String, v std.Int) std.Unit {
			return std.Print(k.ToStr().Plus(std.NewString("=")).Plus(v.ToStr()))
		})
		return std.TheUnit
	}))
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
