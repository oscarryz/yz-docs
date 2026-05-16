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

type Transfer struct {
	std.Cown
	src *Account
	dst *Account
	amount std.Int
}

func NewTransfer(src *Account, dst *Account, amount std.Int) *Transfer {
	return &Transfer{
		src: src,
		dst: dst,
		amount: amount,
	}
}

func (self *Transfer) Run() *std.Thunk[std.Unit] {
	return std.ScheduleMulti([]*std.Cown{&self.Cown, &self.src.Cown, &self.dst.Cown}, func() std.Unit {
		self.src.balance = self.src.balance.Minus(self.amount)
		self.dst.balance = self.dst.balance.Plus(self.amount)
		return std.TheUnit
	})
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.NewThunk(func() std.Unit {
		_bg0 := &std.BocGroup{}
		var alice *Account
		var bob *Account
		std.Schedule(&self.Cown, func() std.Unit {
			alice = NewAccount(std.NewInt(100))
			bob = NewAccount(std.NewInt(0))
			_bg0.GoWait(NewTransfer(alice, bob, std.NewInt(30)).Run())
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		std.Print(alice.balance)
		std.Print(bob.balance)
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
