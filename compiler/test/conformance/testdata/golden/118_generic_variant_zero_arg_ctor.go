package main

import std "yz/runtime/rt"

type _OptionVariant int

const (
	_OptionSome _OptionVariant = iota
	_OptionNone
)

type Option[T any] struct {
	_variant _OptionVariant
	value T
}

func NewOptionSome[T any](value T) *Option[T] {
	return &Option[T]{
		_variant: _OptionSome,
		value: value,
	}
}

func NewOptionNone[T any]() *Option[T] {
	return &Option[T]{
		_variant: _OptionNone,
	}
}

func (self *Option[T]) String() string {
	switch self._variant {
	case _OptionSome:
		return "Option.Some(value: " + std.StringifyRepr(self.value) + ")"
	case _OptionNone:
		return "Option.None()"
	}
	return "Option(?)"
}

func find_first[T any](xs std.Array[T]) *std.Thunk[*Option[T]] {
	return std.Go(func() *Option[T] {
		return func() *Option[T] {
			if xs.Length().Eqeq(std.NewInt(0)).GoBool() {
				return NewOptionNone[T]()
			} else {
				return NewOptionSome(xs.At(std.NewInt(0)))
			}
		}()
	})
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
		var r *Option[std.Int]
		std.Schedule(&self.Cown, func() std.Unit {
			_st0 := find_first(std.NewArray(std.NewInt(1), std.NewInt(2), std.NewInt(3)))
			_bg0.Add(func() { r = _st0.Force() })
			return std.TheUnit
		}).Force()
		_bg0.Wait()
		switch r._variant {
		case _OptionSome:
			std.Print(std.NewString("found: ").Plus(r.value.ToStr()))
		default:
			std.Print(std.NewString("empty"))
		}
		return std.TheUnit
	}))
}

var Main = &_mainBoc{}

func main() {
	Main.Call().Force()
}
