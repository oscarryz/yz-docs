// Package yzrt is the Yz runtime library imported by generated Go code.
// Generated code refers to these types as std.Int, std.String, etc.
// (the "std" alias is set in the generated import block).
package rt

import (
	"fmt"
	"math"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

// ---------------------------------------------------------------------------
// Cycle detection for Stringify
// ---------------------------------------------------------------------------

// stringifyInProgress tracks pointer addresses currently being stringified,
// keyed by (goroutineID, pointer) to be safe under concurrent printing.
var stringifyInProgress sync.Map // key: visitKey, value: struct{}

type visitKey struct {
	gid uint64
	ptr uintptr
}

// goroutineID extracts the current goroutine's ID from the stack trace.
// Used only for cycle detection in Stringify — not a hot path.
func goroutineID() uint64 {
	var buf [32]byte
	n := runtime.Stack(buf[:], false)
	s := string(buf[:n])
	s = strings.TrimPrefix(s, "goroutine ")
	if idx := strings.IndexByte(s, ' '); idx >= 0 {
		id, _ := strconv.ParseUint(s[:idx], 10, 64)
		return id
	}
	return 0
}

// markVisiting registers ptr as in-progress for the current goroutine.
// Returns false if already visited (cycle detected), true if newly registered.
// When true, the caller must call unmarkVisiting(ptr) when done.
func markVisiting(v any) (uintptr, bool) {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() || rv.Kind() != reflect.Ptr || rv.IsNil() {
		return 0, false
	}
	ptr := rv.Pointer()
	key := visitKey{gid: goroutineID(), ptr: ptr}
	_, loaded := stringifyInProgress.LoadOrStore(key, struct{}{})
	if loaded {
		return ptr, false // cycle
	}
	return ptr, true
}

func unmarkVisiting(ptr uintptr) {
	stringifyInProgress.Delete(visitKey{gid: goroutineID(), ptr: ptr})
}

// ---------------------------------------------------------------------------
// Lazy state structs
//
// Each scalar type optionally carries a *lazyX pointer. nil means the value
// is already resolved. Non-nil means the value is backed by a pending
// goroutine; the first access forces it via sync.Once. Because the pointer
// is shared across all copies of a scalar value, forcing through any copy
// resolves all copies simultaneously.
// ---------------------------------------------------------------------------

type lazyInt struct {
	once sync.Once
	val  int64
	fn   func() int64
}

type lazyStr struct {
	once sync.Once
	val  string
	fn   func() string
}

type lazyBool struct {
	once sync.Once
	val  bool
	fn   func() bool
}

type lazyDec struct {
	once sync.Once
	val  float64
	fn   func() float64
}

type lazyUnit struct {
	once sync.Once
	fn   func()
}

// ---------------------------------------------------------------------------
// Int
// ---------------------------------------------------------------------------

// Int is the boxed integer type. It may carry a lazy computation.
type Int struct {
	lazy *lazyInt
	val  int64
}

// NewInt constructs a resolved Int from a Go int64.
func NewInt(v int64) Int { return Int{val: v} }

// LazyInt wraps a *Thunk[Int] into a lazy Int that forces on first access.
func LazyInt(th *Thunk[Int]) Int {
	return Int{lazy: &lazyInt{fn: func() int64 { return th.Force().GoInt() }}}
}

// GoInt returns the underlying Go int64, forcing the lazy computation if needed.
func (i Int) GoInt() int64 {
	if i.lazy == nil {
		return i.val
	}
	i.lazy.once.Do(func() { i.lazy.val = i.lazy.fn() })
	return i.lazy.val
}

// Await blocks until the lazy computation completes.
func (i Int) Await() { i.GoInt() }

func (i Int) String() string { return strconv.FormatInt(i.GoInt(), 10) }

// Arithmetic
func (i Int) Plus(other Int) Int    { return Int{val: i.GoInt() + other.GoInt()} }
func (i Int) Minus(other Int) Int   { return Int{val: i.GoInt() - other.GoInt()} }
func (i Int) Star(other Int) Int    { return Int{val: i.GoInt() * other.GoInt()} }
func (i Int) Slash(other Int) Int   { return Int{val: i.GoInt() / other.GoInt()} }
func (i Int) Percent(other Int) Int { return Int{val: i.GoInt() % other.GoInt()} }

// Unary negation
func (i Int) Neg() Int { return Int{val: -i.GoInt()} }

// Comparison
func (i Int) Lt(other Int) Bool   { return Bool{val: i.GoInt() < other.GoInt()} }
func (i Int) Gt(other Int) Bool   { return Bool{val: i.GoInt() > other.GoInt()} }
func (i Int) Lteq(other Int) Bool { return Bool{val: i.GoInt() <= other.GoInt()} }
func (i Int) Gteq(other Int) Bool { return Bool{val: i.GoInt() >= other.GoInt()} }
func (i Int) Eqeq(other Int) Bool { return Bool{val: i.GoInt() == other.GoInt()} }
func (i Int) Neq(other Int) Bool  { return Bool{val: i.GoInt() != other.GoInt()} }

// Abs returns the absolute value of i.
func (i Int) Abs() Int {
	v := i.GoInt()
	if v < 0 {
		return Int{val: -v}
	}
	return Int{val: v}
}

// To produces a half-open range [i, end).
func (i Int) To(end Int) Range { return Range{from: i.GoInt(), to: end.GoInt()} }

// ToStr converts to String.
func (i Int) ToStr() String { return String{val: strconv.FormatInt(i.GoInt(), 10)} }

// ---------------------------------------------------------------------------
// Decimal
// ---------------------------------------------------------------------------

// Decimal is the boxed floating-point type. It may carry a lazy computation.
type Decimal struct {
	lazy *lazyDec
	val  float64
}

// NewDecimal constructs a resolved Decimal from a Go float64.
func NewDecimal(v float64) Decimal { return Decimal{val: v} }

// LazyDecimal wraps a *Thunk[Decimal] into a lazy Decimal that forces on first access.
func LazyDecimal(th *Thunk[Decimal]) Decimal {
	return Decimal{lazy: &lazyDec{fn: func() float64 { return th.Force().GoFloat64() }}}
}

// GoFloat64 returns the underlying Go float64, forcing the lazy computation if needed.
func (d Decimal) GoFloat64() float64 {
	if d.lazy == nil {
		return d.val
	}
	d.lazy.once.Do(func() { d.lazy.val = d.lazy.fn() })
	return d.lazy.val
}

// Await blocks until the lazy computation completes.
func (d Decimal) Await() { d.GoFloat64() }

func (d Decimal) String() string { return strconv.FormatFloat(d.GoFloat64(), 'g', -1, 64) }

// Arithmetic
func (d Decimal) Plus(other Decimal) Decimal  { return Decimal{val: d.GoFloat64() + other.GoFloat64()} }
func (d Decimal) Minus(other Decimal) Decimal { return Decimal{val: d.GoFloat64() - other.GoFloat64()} }
func (d Decimal) Star(other Decimal) Decimal  { return Decimal{val: d.GoFloat64() * other.GoFloat64()} }
func (d Decimal) Slash(other Decimal) Decimal { return Decimal{val: d.GoFloat64() / other.GoFloat64()} }

// Unary negation
func (d Decimal) Neg() Decimal { return Decimal{val: -d.GoFloat64()} }

// Abs returns the absolute value of d.
func (d Decimal) Abs() Decimal { return Decimal{val: math.Abs(d.GoFloat64())} }

// Pow raises d to the power of exp.
func (d Decimal) Pow(exp Decimal) Decimal {
	return Decimal{val: math.Pow(d.GoFloat64(), exp.GoFloat64())}
}

// Comparison
func (d Decimal) Lt(other Decimal) Bool   { return Bool{val: d.GoFloat64() < other.GoFloat64()} }
func (d Decimal) Gt(other Decimal) Bool   { return Bool{val: d.GoFloat64() > other.GoFloat64()} }
func (d Decimal) Lteq(other Decimal) Bool { return Bool{val: d.GoFloat64() <= other.GoFloat64()} }
func (d Decimal) Gteq(other Decimal) Bool { return Bool{val: d.GoFloat64() >= other.GoFloat64()} }
func (d Decimal) Eqeq(other Decimal) Bool { return Bool{val: d.GoFloat64() == other.GoFloat64()} }
func (d Decimal) Neq(other Decimal) Bool  { return Bool{val: d.GoFloat64() != other.GoFloat64()} }

// ToStr converts to String.
func (d Decimal) ToStr() String {
	return String{val: strconv.FormatFloat(d.GoFloat64(), 'g', -1, 64)}
}

// ---------------------------------------------------------------------------
// String
// ---------------------------------------------------------------------------

// String is the boxed string type. It may carry a lazy computation.
type String struct {
	lazy *lazyStr
	val  string
}

// NewString constructs a resolved String from a Go string.
func NewString(v string) String { return String{val: v} }

// LazyString wraps a *Thunk[String] into a lazy String that forces on first access.
func LazyString(th *Thunk[String]) String {
	return String{lazy: &lazyStr{fn: func() string { return th.Force().GoString() }}}
}

// GoString returns the underlying Go string, forcing the lazy computation if needed.
func (s String) GoString() string {
	if s.lazy == nil {
		return s.val
	}
	s.lazy.once.Do(func() { s.lazy.val = s.lazy.fn() })
	return s.lazy.val
}

// Await blocks until the lazy computation completes.
func (s String) Await() { s.GoString() }

func (s String) String() string { return s.GoString() }

// Plus concatenates two strings.
func (s String) Plus(other String) String { return String{val: s.GoString() + other.GoString()} }

// Comparison
func (s String) Eqeq(other String) Bool { return Bool{val: s.GoString() == other.GoString()} }
func (s String) Neq(other String) Bool  { return Bool{val: s.GoString() != other.GoString()} }
func (s String) Lt(other String) Bool   { return Bool{val: s.GoString() < other.GoString()} }
func (s String) Gt(other String) Bool   { return Bool{val: s.GoString() > other.GoString()} }
func (s String) Lteq(other String) Bool { return Bool{val: s.GoString() <= other.GoString()} }
func (s String) Gteq(other String) Bool { return Bool{val: s.GoString() >= other.GoString()} }

// Length returns the number of Unicode code points.
func (s String) Length() Int { return Int{val: int64(len([]rune(s.GoString())))} }

// Contains reports whether sub is within s.
func (s String) Contains(sub String) Bool {
	return Bool{val: strings.Contains(s.GoString(), sub.GoString())}
}

// HasPrefix reports whether s begins with prefix.
func (s String) HasPrefix(prefix String) Bool {
	return Bool{val: strings.HasPrefix(s.GoString(), prefix.GoString())}
}

// HasSuffix reports whether s ends with suffix.
func (s String) HasSuffix(suffix String) Bool {
	return Bool{val: strings.HasSuffix(s.GoString(), suffix.GoString())}
}

// ToUpper returns s with all letters mapped to upper case.
func (s String) ToUpper() String {
	return String{val: strings.Map(unicode.ToUpper, s.GoString())}
}

// ToLower returns s with all letters mapped to lower case.
func (s String) ToLower() String {
	return String{val: strings.Map(unicode.ToLower, s.GoString())}
}

// Trim returns s with leading and trailing white space removed.
func (s String) Trim() String { return String{val: strings.TrimSpace(s.GoString())} }

// ToStr returns the string itself.
func (s String) ToStr() String { return s }

// ---------------------------------------------------------------------------
// Bool
// ---------------------------------------------------------------------------

// Bool is the boxed boolean type. It may carry a lazy computation.
type Bool struct {
	lazy *lazyBool
	val  bool
}

// NewBool constructs a resolved Bool from a Go bool.
func NewBool(v bool) Bool { return Bool{val: v} }

// LazyBool wraps a *Thunk[Bool] into a lazy Bool that forces on first access.
func LazyBool(th *Thunk[Bool]) Bool {
	return Bool{lazy: &lazyBool{fn: func() bool { return th.Force().GoBool() }}}
}

// GoBool returns the underlying Go bool, forcing the lazy computation if needed.
func (b Bool) GoBool() bool {
	if b.lazy == nil {
		return b.val
	}
	b.lazy.once.Do(func() { b.lazy.val = b.lazy.fn() })
	return b.lazy.val
}

// Await blocks until the lazy computation completes.
func (b Bool) Await() { b.GoBool() }

func (b Bool) String() string {
	if b.GoBool() {
		return "true"
	}
	return "false"
}

// Logical
func (b Bool) Ampamp(other Bool) Bool   { return Bool{val: b.GoBool() && other.GoBool()} }
func (b Bool) Pipepipe(other Bool) Bool { return Bool{val: b.GoBool() || other.GoBool()} }

// Comparison
func (b Bool) Eqeq(other Bool) Bool { return Bool{val: b.GoBool() == other.GoBool()} }
func (b Bool) Neq(other Bool) Bool  { return Bool{val: b.GoBool() != other.GoBool()} }

// Qm is the conditional operator: flag ? { trueCase } , { falseCase }
// Both branches are passed as zero-argument functions so only the selected
// branch is evaluated.
func (b Bool) Qm(trueCase, falseCase func() any) any {
	if b.GoBool() {
		return trueCase()
	}
	return falseCase()
}

// ToStr converts to String.
func (b Bool) ToStr() String { return NewString(b.String()) }

// ---------------------------------------------------------------------------
// Unit
// ---------------------------------------------------------------------------

// Unit is the empty/void type — the result of a boc that returns nothing.
// It may carry a lazy computation (from a boc call that returns unit).
type Unit struct {
	lazy *lazyUnit
}

// TheUnit is the singleton resolved Unit value.
var TheUnit = Unit{}

// LazyUnit wraps a *Thunk[Unit] into a lazy Unit that forces on Await.
func LazyUnit(th *Thunk[Unit]) Unit {
	return Unit{lazy: &lazyUnit{fn: func() { th.Force() }}}
}

// Await blocks until the lazy computation completes.
func (u Unit) Await() {
	if u.lazy != nil {
		u.lazy.once.Do(u.lazy.fn)
	}
}

// Force is an alias for Await, kept for the entry-point call Main.Call().Force().
func (u Unit) Force() { u.Await() }

func (u Unit) String() string { return "()" }

// ---------------------------------------------------------------------------
// Stringer helper (for Print / Info)
// ---------------------------------------------------------------------------

// Stringify returns a human-readable string for any yzrt value.
// It handles all boxed types plus raw Go values via fmt.Sprint.
func Stringify(v any) string {
	if v == nil {
		return "_"
	}
	if rv := reflect.ValueOf(v); rv.Kind() == reflect.Ptr && rv.IsNil() {
		return "_"
	}
	switch x := v.(type) {
	case fmt.Stringer:
		ptr, ok := markVisiting(v)
		if !ok && ptr != 0 {
			return YzTypeName(v) + "(...)"
		}
		if ok {
			defer unmarkVisiting(ptr)
		}
		return x.String()
	default:
		return fmt.Sprint(v)
	}
}

// YzTypeName returns the Yz type name for a value (used in homoiconic output).
func YzTypeName(v any) string {
	switch v.(type) {
	case String:
		return "String"
	case Int:
		return "Int"
	case Bool:
		return "Bool"
	case Decimal:
		return "Decimal"
	default:
		t := reflect.TypeOf(v)
		if t == nil {
			return "?"
		}
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		return t.Name()
	}
}

// StringifyRepr returns a Yz-syntax representation of v.
// Unlike Stringify, String values are quoted so they are distinguishable
// from identifiers inside homoiconic output (backtick interpolation).
// Cycles are detected per-goroutine; a cyclic reference prints as TypeName(...).
func StringifyRepr(v any) string {
	if v == nil {
		return "_"
	}
	if s, ok := v.(String); ok {
		return `"` + s.GoString() + `"`
	}
	// Nil pointer wrapped in interface: don't call String() on it.
	if rv := reflect.ValueOf(v); rv.Kind() == reflect.Ptr && rv.IsNil() {
		return "_"
	}
	switch x := v.(type) {
	case fmt.Stringer:
		ptr, ok := markVisiting(v)
		if !ok && ptr != 0 {
			return YzTypeName(v) + "(...)"
		}
		if ok {
			defer unmarkVisiting(ptr)
		}
		return x.String()
	default:
		return fmt.Sprint(v)
	}
}
