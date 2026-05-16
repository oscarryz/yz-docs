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

type _loaderBoc struct {
	std.Cown
	acc *Account
}

func (self *_loaderBoc) Call(acc *Account) *std.Thunk[*Account] {
	return std.ScheduleMulti([]*std.Cown{&self.Cown, &acc.Cown}, func() *Account {
		self.acc = acc
		return self.acc
	})
}

var Loader = &_loaderBoc{
}

type _userBoc struct {
	std.Cown
	acc *Account
}

func (self *_userBoc) Call(acc *Account) *std.Thunk[std.Unit] {
	return std.ScheduleMulti([]*std.Cown{&self.Cown, &acc.Cown}, func() std.Unit {
		self.acc = acc
		var loaded *Account
		return std.Print(loaded.balance)
	})
}

var User = &_userBoc{
}

type _mainBoc struct {
	std.Cown
}

func (self *_mainBoc) Call() *std.Thunk[std.Unit] {
	return std.NewThunk(func() std.Unit {
		_bg0 := &std.BocGroup{}
		var a *Account
		std.Schedule(&self.Cown, func() std.Unit {
			a = NewAccount(std.NewInt(42))
			_bg0.GoWait((&_userBoc{}).Call(a))
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		return std.TheUnit
	})
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
