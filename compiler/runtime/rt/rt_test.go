package rt

import (
	"sync"
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

// ---------------------------------------------------------------------------
// Cown / Schedule
// ---------------------------------------------------------------------------

func TestScheduleBasic(t *testing.T) {
	var c Cown
	th := Schedule(&c, func() Int { return NewInt(42) })
	if th.Force().GoInt() != 42 {
		t.Error("Schedule basic")
	}
}

// TestScheduleSerializes verifies that N concurrent Schedule calls on the same
// cown serialize: the counter ends up at N with no lost updates.
func TestScheduleSerializes(t *testing.T) {
	const n = 1000
	var c Cown
	count := 0 // unprotected — Schedule must serialize access

	var wg sync.WaitGroup
	wg.Add(n)
	for range n {
		go func() {
			Schedule(&c, func() Unit {
				count++
				return TheUnit
			}).Force()
			wg.Done()
		}()
	}
	wg.Wait()

	if count != n {
		t.Errorf("lost updates: got %d, want %d", count, n)
	}
}

// TestSchedulePreservesOrder verifies that behaviours on the same cown run in
// the order they were spawned (FIFO).
func TestSchedulePreservesOrder(t *testing.T) {
	const n = 100
	var c Cown
	var order []int

	var wg sync.WaitGroup
	wg.Add(n)
	// Serialize the spawning so positions are well-defined.
	for i := range n {
		i := i
		Schedule(&c, func() Unit {
			order = append(order, i)
			wg.Done()
			return TheUnit
		}).Force()
	}
	wg.Wait()

	for i, v := range order {
		if v != i {
			t.Fatalf("out of order at position %d: got %d", i, v)
		}
	}
}

// TestScheduleTwoIndependentCowns verifies that behaviours on different cowns
// run in parallel (both can make progress without waiting for each other).
func TestScheduleTwoIndependentCowns(t *testing.T) {
	var c1, c2 Cown
	ch1 := make(chan struct{})
	ch2 := make(chan struct{})

	t1 := Schedule(&c1, func() Unit {
		close(ch1)
		<-ch2 // wait for c2's behaviour — would deadlock if serialized
		return TheUnit
	})
	t2 := Schedule(&c2, func() Unit {
		close(ch2)
		<-ch1
		return TheUnit
	})

	t1.Force()
	t2.Force()
}

