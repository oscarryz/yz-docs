package rt

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

// Filter returns a new Array containing only elements for which fn returns Bool true.
func (a Array[T]) Filter(fn func(T) Bool) Array[T] {
	var result []T
	for _, v := range a.elems {
		if fn(v).val {
			result = append(result, v)
		}
	}
	return Array[T]{elems: result}
}

// Each calls fn for every element in the array.
func (a Array[T]) Each(fn func(T) Unit) {
	for _, v := range a.elems {
		fn(v)
	}
}

// Any reports whether fn returns true for at least one element.
func (a Array[T]) Any(fn func(T) Bool) Bool {
	for _, v := range a.elems {
		if fn(v).val {
			return Bool{true}
		}
	}
	return Bool{false}
}

// All reports whether fn returns true for every element.
func (a Array[T]) All(fn func(T) Bool) Bool {
	for _, v := range a.elems {
		if !fn(v).val {
			return Bool{false}
		}
	}
	return Bool{true}
}

// IsEmpty reports whether the array has no elements.
func (a Array[T]) IsEmpty() Bool { return Bool{len(a.elems) == 0} }

// ArrayMap applies fn to each element of a and returns a new Array of results.
// It is a package-level function because Go methods cannot introduce new type parameters.
func ArrayMap[T, U any](a Array[T], fn func(T) U) Array[U] {
	result := make([]U, len(a.elems))
	for i, v := range a.elems {
		result[i] = fn(v)
	}
	return Array[U]{elems: result}
}

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
