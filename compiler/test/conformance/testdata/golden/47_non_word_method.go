package main

import std "yz/runtime/rt"

type Greeter struct {
	std.Cown
	name std.String
}

func NewGreeter(name std.String) *Greeter {
	return &Greeter{
		name: name,
	}
}

func (self *Greeter) String() string {
	return "Greeter(name: " + std.StringifyRepr(self.name) + ")"
}

func (self *Greeter) plusplus(other std.String) std.String {
	return self.name.ToStr().Plus(std.NewString(" and ")).Plus(other.ToStr())
}

func (self *Greeter) Plusplus(other std.String) std.String {
	return std.LazyString(std.Schedule(&self.Cown, func() std.String {
		return self.plusplus(other)
	}))
}

func (self *Greeter) Name() std.String {
	return self.name
}

type Tag struct {
	std.Cown
	label std.String
}

func NewTag(label std.String) *Tag {
	return &Tag{
		label: label,
	}
}

func (self *Tag) String() string {
	return "Tag(label: " + std.StringifyRepr(self.label) + ")"
}

func (self *Tag) plusplus(other *Tag) *Tag {
	return NewTag(std.NewString("(").Plus(self.label.ToStr()).Plus(std.NewString(" ")).Plus(other.label.ToStr()).Plus(std.NewString(")")))
}

func (self *Tag) Plusplus(other *Tag) *std.Thunk[*Tag] {
	return std.Schedule(&self.Cown, func() *Tag {
		return self.plusplus(other)
	})
}

func (self *Tag) Label() std.String {
	return self.label
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
		var c std.String
		var a *Greeter
		std.Schedule(&self.Cown, func() std.Unit {
			a = NewGreeter(std.NewString("Alice"))
			c = a.Plusplus(std.NewString("Bob"))
			_bg0.Add(func() { c.Await() })
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		std.Print(c)
		var x *Tag = NewTag(std.NewString("a"))
		var y *Tag = NewTag(std.NewString("b"))
		var z *Tag = NewTag(std.NewString("c"))
		_bg1 := &std.BocGroup{}
		var left *Tag
		_th0 := x.Plusplus(y).Force().Plusplus(z)
		_bg1.Add(func() { left = _th0.Force() })
		_bg1.Wait()
		std.Print(left.label)
		_bg2 := &std.BocGroup{}
		var yz *Tag
		_th1 := y.Plusplus(z)
		_bg2.Add(func() { yz = _th1.Force() })
		_bg2.Wait()
		_bg3 := &std.BocGroup{}
		var right *Tag
		_th2 := x.Plusplus(yz)
		_bg3.Add(func() { right = _th2.Force() })
		_bg3.Wait()
		std.Print(right.label)
		return std.TheUnit
	}))
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
