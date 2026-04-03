package yzrt

import "fmt"

// ---------------------------------------------------------------------------
// Array[T]
// ---------------------------------------------------------------------------

// Array is the ordered, homogeneous collection type.
type Array[T any] struct {
	elems []T
}

// NewArray constructs an Array from a Go slice.
func NewArray[T any](elems ...T) Array[T] {
	cp := make([]T, len(elems))
	copy(cp, elems)
	return Array[T]{elems: cp}
}

// At returns the element at index i (panics on out-of-bounds, consistent with Yz spec).
func (a Array[T]) At(i Int) T {
	return a.elems[i.val]
}

// Set returns a new Array with element at index i replaced by v.
func (a Array[T]) Set(i Int, v T) Array[T] {
	cp := make([]T, len(a.elems))
	copy(cp, a.elems)
	cp[i.val] = v
	return Array[T]{elems: cp}
}

// Append returns a new Array with v appended.
func (a Array[T]) Append(v T) Array[T] {
	cp := make([]T, len(a.elems)+1)
	copy(cp, a.elems)
	cp[len(a.elems)] = v
	return Array[T]{elems: cp}
}

// Length returns the number of elements.
func (a Array[T]) Length() Int { return Int{int64(len(a.elems))} }

// GoSlice returns the underlying Go slice (for interop / codegen helpers).
func (a Array[T]) GoSlice() []T { return a.elems }

func (a Array[T]) String() string { return fmt.Sprintf("%v", a.elems) }

// ---------------------------------------------------------------------------
// Dict[K, V]
// ---------------------------------------------------------------------------

// Dict is the unordered key-value collection type.
// K must be comparable.
type Dict[K comparable, V any] struct {
	m map[K]V
}

// NewDict constructs an empty Dict.
func NewDict[K comparable, V any]() Dict[K, V] {
	return Dict[K, V]{m: make(map[K]V)}
}

// At returns the value for key k. Panics if key not present.
func (d Dict[K, V]) At(k K) V {
	v, ok := d.m[k]
	if !ok {
		panic(fmt.Sprintf("dict: key not found: %v", k))
	}
	return v
}

// Set returns a new Dict with k mapped to v.
func (d Dict[K, V]) Set(k K, v V) Dict[K, V] {
	m := make(map[K]V, len(d.m)+1)
	for kk, vv := range d.m {
		m[kk] = vv
	}
	m[k] = v
	return Dict[K, V]{m: m}
}

// Has reports whether key k exists.
func (d Dict[K, V]) Has(k K) Bool {
	_, ok := d.m[k]
	return Bool{ok}
}

// Length returns the number of key-value pairs.
func (d Dict[K, V]) Length() Int { return Int{int64(len(d.m))} }

// GoMap returns the underlying Go map (for interop / codegen helpers).
func (d Dict[K, V]) GoMap() map[K]V { return d.m }

func (d Dict[K, V]) String() string { return fmt.Sprintf("%v", d.m) }

// ---------------------------------------------------------------------------
// Range
// ---------------------------------------------------------------------------

// Range is a half-open integer interval [from, to) produced by Int.To().
type Range struct {
	from, to int64
}

// NewRange constructs a Range.
func NewRange(from, to int64) Range { return Range{from: from, to: to} }

// Each calls fn for every integer in the range.
func (r Range) Each(fn func(Int)) {
	for i := r.from; i < r.to; i++ {
		fn(Int{i})
	}
}

// ToArray materializes the range into an Array[Int].
func (r Range) ToArray() Array[Int] {
	size := r.to - r.from
	if size <= 0 {
		return Array[Int]{}
	}
	elems := make([]Int, size)
	for i := int64(0); i < size; i++ {
		elems[i] = Int{r.from + i}
	}
	return Array[Int]{elems: elems}
}

// Length returns the number of integers in the range.
func (r Range) Length() Int {
	if r.to <= r.from {
		return Int{0}
	}
	return Int{r.to - r.from}
}

func (r Range) String() string { return fmt.Sprintf("%d..%d", r.from, r.to) }
