package rt

import "sync"

// Thunk[T] is a lazy value produced by every boc invocation.
// The computation runs in a goroutine; the result is materialized
// on the first call to Force(), which blocks until the goroutine completes.
type Thunk[T any] struct {
	once sync.Once
	val  T
	fn   func() T
}

// NewThunk creates a Thunk that will evaluate fn (in the caller's goroutine
// or in a spawned goroutine) and cache the result.
func NewThunk[T any](fn func() T) *Thunk[T] {
	return &Thunk[T]{fn: fn}
}

// Force materializes the thunk, blocking until the value is available.
// Subsequent calls return the cached value immediately.
func (th *Thunk[T]) Force() T {
	th.once.Do(func() { th.val = th.fn() })
	return th.val
}

// Go launches fn in a new goroutine and returns a Thunk that materializes
// the result. This is the standard way all boc calls are compiled:
//
//	result := yzrt.Go(func() Int { return a.plus(b) })
//	// ... later:
//	val := result.Force()
func Go[T any](fn func() T) *Thunk[T] {
	th := &Thunk[T]{fn: fn}
	// Channel used to signal completion so Force() can sync.Once safely.
	done := make(chan struct{})
	go func() {
		th.val = fn()
		close(done)
	}()
	// Replace fn with one that waits for the goroutine then returns cached val.
	th.fn = func() T {
		<-done
		return th.val
	}
	return th
}
