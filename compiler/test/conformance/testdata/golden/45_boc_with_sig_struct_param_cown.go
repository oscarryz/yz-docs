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

func (self *Account) Balance() std.Int {
	return self.balance
}

type _transferBoc struct {
	std.Cown
	src *Account
	dst *Account
	amount std.Int
}

func (self *_transferBoc) String() string {
	return "{ " + "src: " + std.StringifyRepr(self.src) + "; " + "dst: " + std.StringifyRepr(self.dst) + "; " + "amount: " + std.StringifyRepr(self.amount) + "; " + "call: {}" + " }"
}

func (self *_transferBoc) Call(src *Account, dst *Account, amount std.Int) std.Unit {
	return std.LazyUnit(std.ScheduleMulti([]*std.Cown{&self.Cown, &src.Cown, &dst.Cown}, func() std.Unit {
		self.src = src
		self.dst = dst
		self.amount = amount
		self.src.balance = self.src.balance.Minus(self.amount)
		self.dst.balance = self.dst.balance.Plus(self.amount)
		return std.TheUnit
	}))
}

var Transfer = &_transferBoc{
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
		var alice *Account
		var bob *Account
		std.Schedule(&self.Cown, func() std.Unit {
			alice = NewAccount(std.NewInt(100))
			bob = NewAccount(std.NewInt(0))
			_st0 := Transfer.Call(alice, bob, std.NewInt(30))
			_bg0.Add(func() { _st0.Await() })
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		std.Print(std.NewString(std.StringifyRepr(alice.balance)))
		std.Print(std.NewString(std.StringifyRepr(bob.balance)))
		return std.TheUnit
	}))
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
