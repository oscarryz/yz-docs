package main

import std "yz/runtime/rt"

type _bankBoc struct {
	std.Cown
	balance std.Int
}

func (self *_bankBoc) Deposit(amount std.Int) *std.Thunk[std.Unit] {
	return std.Schedule(&self.Cown, func() std.Unit {
		self.balance = self.balance.Plus(amount)
		return std.TheUnit
	})
}

var Bank = &_bankBoc{
	balance: std.NewInt(0),
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.Schedule(&self.Cown, func() std.Unit {
		std.Schedule(&Bank.Cown, func() std.Unit {
			Bank.balance = std.NewInt(42)
			return std.TheUnit
		}).Force()
		std.Print(Bank.balance)
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
