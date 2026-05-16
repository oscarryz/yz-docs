package main

import std "yz/runtime/rt"

type _bankBoc struct {
	std.Cown
	balance std.Int
}

func (self *_bankBoc) deposit(amount std.Int) std.Unit {
	self.balance = self.balance.Plus(amount)
	return std.TheUnit
}

func (self *_bankBoc) Deposit(amount std.Int) *std.Thunk[std.Unit] {
	return std.Schedule(&self.Cown, func() std.Unit {
		return self.deposit(amount)
	})
}

var Bank = &_bankBoc{
	balance: std.NewInt(0),
}

type _ledgerBoc struct {
	std.Cown
	total std.Int
}

func (self *_ledgerBoc) add(amount std.Int) std.Unit {
	self.total = self.total.Plus(amount)
	return std.TheUnit
}

func (self *_ledgerBoc) Add(amount std.Int) *std.Thunk[std.Unit] {
	return std.Schedule(&self.Cown, func() std.Unit {
		return self.add(amount)
	})
}

var Ledger = &_ledgerBoc{
	total: std.NewInt(0),
}

type _syncBoc struct {
	std.Cown
	b *_bankBoc
	l *_ledgerBoc
}

func (self *_syncBoc) Call(b *_bankBoc, l *_ledgerBoc) *std.Thunk[std.Unit] {
	return std.ScheduleMulti([]*std.Cown{&self.Cown, &b.Cown, &l.Cown}, func() std.Unit {
		self.b = b
		self.l = l
		self.b.balance = self.b.balance.Plus(std.NewInt(1))
		self.l.total = self.l.total.Plus(std.NewInt(1))
		return std.TheUnit
	})
}

var Sync = &_syncBoc{
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.NewThunk(func() std.Unit {
		_bg0 := &std.BocGroup{}
		std.Schedule(&self.Cown, func() std.Unit {
			_bg0.GoWait((&_syncBoc{}).Call(Bank, Ledger))
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		std.Print(Bank.balance)
		std.Print(Ledger.total)
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
