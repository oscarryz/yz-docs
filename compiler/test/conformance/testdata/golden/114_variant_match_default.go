package main

import std "yz/runtime/rt"

type _PetVariant int

const (
	_PetCat _PetVariant = iota
	_PetDog
)

type Pet struct {
	_variant _PetVariant
	name std.String
}

func NewPetCat(name std.String) *Pet {
	return &Pet{
		_variant: _PetCat,
		name: name,
	}
}

func NewPetDog(name std.String) *Pet {
	return &Pet{
		_variant: _PetDog,
		name: name,
	}
}

func (self *Pet) String() string {
	switch self._variant {
	case _PetCat:
		return "Pet.Cat(name: " + std.StringifyRepr(self.name) + ")"
	case _PetDog:
		return "Pet.Dog(name: " + std.StringifyRepr(self.name) + ")"
	}
	return "Pet(?)"
}

type _describeBoc struct {
	std.Cown
	pet *Pet
}

func (self *_describeBoc) String() string {
	return "{ " + "pet: " + std.StringifyRepr(self.pet) + "; " + "call: {}" + " }"
}

func (self *_describeBoc) Call(pet *Pet) std.String {
	return std.LazyString(std.Schedule(&self.Cown, func() std.String {
		self.pet = pet
		return func() std.String {
			switch self.pet._variant {
			case _PetCat:
				return std.NewString("cat ").Plus(self.pet.name.ToStr())
			default:
				return std.NewString("unknown")
			}
			return std.NewString("")
		}()
	}))
}

var Describe = &_describeBoc{
}

type _announceBoc struct {
	std.Cown
	pet *Pet
}

func (self *_announceBoc) String() string {
	return "{ " + "pet: " + std.StringifyRepr(self.pet) + "; " + "call: {}" + " }"
}

func (self *_announceBoc) Call(pet *Pet) std.Unit {
	return std.LazyUnit(std.Schedule(&self.Cown, func() std.Unit {
		self.pet = pet
		return func() std.Unit {
			switch self.pet._variant {
			case _PetCat:
				return std.Print(std.NewString("cat: ").Plus(self.pet.name.ToStr()))
			default:
				return std.Print(std.NewString("other"))
			}
			return std.TheUnit
		}()
	}))
}

var Announce = &_announceBoc{
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
		std.Schedule(&self.Cown, func() std.Unit {
			std.Print(Describe.Call(NewPetCat(std.NewString("Whiskers"))))
			std.Print(Describe.Call(NewPetDog(std.NewString("Rex"))))
			_st0 := Announce.Call(NewPetCat(std.NewString("Whiskers")))
			_bg0.Add(func() { _st0.Await() })
			_st1 := Announce.Call(NewPetDog(std.NewString("Rex")))
			_bg0.Add(func() { _st1.Await() })
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
