package rt

import (
	"testing"
)

// ---------------------------------------------------------------------------
// Int
// ---------------------------------------------------------------------------

func TestIntArithmetic(t *testing.T) {
	a := NewInt(10)
	b := NewInt(3)
	if a.Plus(b).GoInt() != 13 {
		t.Error("Plus")
	}
	if a.Minus(b).GoInt() != 7 {
		t.Error("Minus")
	}
	if a.Star(b).GoInt() != 30 {
		t.Error("Star")
	}
	if a.Slash(b).GoInt() != 3 {
		t.Error("Slash")
	}
	if a.Percent(b).GoInt() != 1 {
		t.Error("Percent")
	}
}

func TestIntComparison(t *testing.T) {
	a, b := NewInt(5), NewInt(10)
	if !a.Lt(b).GoBool() {
		t.Error("Lt")
	}
	if !b.Gt(a).GoBool() {
		t.Error("Gt")
	}
	if !a.Eqeq(NewInt(5)).GoBool() {
		t.Error("Eqeq")
	}
	if !a.Neq(b).GoBool() {
		t.Error("Neq")
	}
	if !a.Lteq(a).GoBool() {
		t.Error("Lteq")
	}
	if !b.Gteq(b).GoBool() {
		t.Error("Gteq")
	}
}

func TestIntRange(t *testing.T) {
	r := NewInt(1).To(NewInt(5))
	if r.Length().GoInt() != 4 {
		t.Errorf("Range length: got %d, want 4", r.Length().GoInt())
	}
	arr := r.ToArray()
	if arr.Length().GoInt() != 4 {
		t.Error("ToArray length")
	}
	if arr.At(NewInt(0)).GoInt() != 1 {
		t.Error("ToArray[0]")
	}
}

func TestIntToStr(t *testing.T) {
	if NewInt(42).ToStr().GoString() != "42" {
		t.Error("ToStr")
	}
}

// ---------------------------------------------------------------------------
// Decimal
// ---------------------------------------------------------------------------

func TestDecimalArithmetic(t *testing.T) {
	a := NewDecimal(1.5)
	b := NewDecimal(0.5)
	if a.Plus(b).GoFloat64() != 2.0 {
		t.Error("Plus")
	}
	if a.Minus(b).GoFloat64() != 1.0 {
		t.Error("Minus")
	}
}

// ---------------------------------------------------------------------------
// String
// ---------------------------------------------------------------------------

func TestStringPlus(t *testing.T) {
	s := NewString("hello").Plus(NewString(" world"))
	if s.GoString() != "hello world" {
		t.Errorf("got %q", s.GoString())
	}
}

func TestStringLength(t *testing.T) {
	if NewString("abc").Length().GoInt() != 3 {
		t.Error("length")
	}
}

// ---------------------------------------------------------------------------
// Bool
// ---------------------------------------------------------------------------

func TestBoolQm(t *testing.T) {
	result := NewBool(true).Qm(
		func() any { return NewString("yes") },
		func() any { return NewString("no") },
	)
	if s, ok := result.(String); !ok || s.GoString() != "yes" {
		t.Errorf("Qm true: got %v", result)
	}

	result = NewBool(false).Qm(
		func() any { return NewString("yes") },
		func() any { return NewString("no") },
	)
	if s, ok := result.(String); !ok || s.GoString() != "no" {
		t.Errorf("Qm false: got %v", result)
	}
}

// ---------------------------------------------------------------------------
// Thunk
// ---------------------------------------------------------------------------

func TestThunkForce(t *testing.T) {
	th := NewThunk(func() Int { return NewInt(42) })
	if th.Force().GoInt() != 42 {
		t.Error("Force")
	}
	// Second call returns cached value.
	if th.Force().GoInt() != 42 {
		t.Error("Force (cached)")
	}
}

func TestThunkGo(t *testing.T) {
	th := Go(func() Int { return NewInt(99) })
	if th.Force().GoInt() != 99 {
		t.Error("Go thunk")
	}
}

// ---------------------------------------------------------------------------
// Array
// ---------------------------------------------------------------------------

func TestArray(t *testing.T) {
	a := NewArray(NewInt(1), NewInt(2), NewInt(3))
	if a.Length().GoInt() != 3 {
		t.Error("Length")
	}
	if a.At(NewInt(1)).GoInt() != 2 {
		t.Error("At")
	}
	a2 := a.Append(NewInt(4))
	if a2.Length().GoInt() != 4 {
		t.Error("Append length")
	}
}

// ---------------------------------------------------------------------------
// Dict
// ---------------------------------------------------------------------------

func TestDict(t *testing.T) {
	d := NewDict[String, Int]()
	d = d.Set(NewString("a"), NewInt(1))
	d = d.Set(NewString("b"), NewInt(2))
	if d.Length().GoInt() != 2 {
		t.Error("Length")
	}
	if d.At(NewString("a")).GoInt() != 1 {
		t.Error("At")
	}
	if !d.Has(NewString("b")).GoBool() {
		t.Error("Has")
	}
}

