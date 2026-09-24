package main

import std "yz/runtime/rt"

type Account struct {
	std.Cown
	balance std.Int
}

func NewAccount(balance std.Int) *Account {
	return &Account{
		balance: balance,
	}
}

func (self *Account) String() string {
	return "Account(balance: " + std.StringifyRepr(self.balance) + ")"
}

func (self *Account) toStr() std.String {
	return self.balance.ToStr()
}

func (self *Account) ToStr() std.String {
	return std.LazyString(std.Schedule(&self.Cown, func() std.String {
		return self.toStr()
	}))
}

func (self *Account) Balance() std.Int {
	return self.balance
}

type _showBoc struct {
	std.Cown
	src *Account
	dst *Account
}

func (self *_showBoc) String() string {
	return "{ " + "src: " + std.StringifyRepr(self.src) + "; " + "dst: " + std.StringifyRepr(self.dst) + "; " + "call: {}" + " }"
}

func (self *_showBoc) Call(src *Account, dst *Account) std.Unit {
	return std.LazyUnit(std.ScheduleMulti([]*std.Cown{&self.Cown, &src.Cown, &dst.Cown}, func() std.Unit {
		self.src = src
		self.dst = dst
		return std.Print(self.src.toStr())
	}))
}

var Show = &_showBoc{
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
		var a *Account
		var b *Account
		std.Schedule(&self.Cown, func() std.Unit {
			a = NewAccount(std.NewInt(10))
			b = NewAccount(std.NewInt(0))
			_st0 := Show.Call(a, b)
			_bg0.Add(func() { _st0.Await() })
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		return std.TheUnit
	}))
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
