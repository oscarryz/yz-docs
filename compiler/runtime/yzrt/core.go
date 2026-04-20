package yzrt

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// ---------------------------------------------------------------------------
// Print
// ---------------------------------------------------------------------------

// Print writes the string representation of v to stdout followed by a newline.
// It accepts any yzrt boxed value or a raw Go value and uses Stringify to
// produce the output.
func Print(v any) Unit {
	fmt.Fprintln(os.Stdout, Stringify(v))
	return TheUnit
}

// ---------------------------------------------------------------------------
// Info
// ---------------------------------------------------------------------------

// InfoResult is returned by Info — it carries the value and provides
// a human-readable representation useful for debugging.
type InfoResult struct {
	Value any
}

func (r InfoResult) String() string {
	b, err := json.Marshal(r.Value)
	if err != nil {
		return fmt.Sprintf("%v", r.Value)
	}
	return string(b)
}

// Info returns an InfoResult wrapping v. It is the runtime of the
// `info(expr)` built-in. The caller may force and print the result.
func Info(v any) InfoResult {
	return InfoResult{Value: v}
}

// ---------------------------------------------------------------------------
// WaitGroup helper for structured concurrency
// ---------------------------------------------------------------------------

// BocGroup manages structured concurrency for a single boc invocation.
// Each spawned child boc registers itself; the parent waits at the end.
type BocGroup struct {
	wg sync.WaitGroup
}

// Go spawns fn as a goroutine registered with this group.
// Returns a Thunk that materializes once fn completes.
func (g *BocGroup) Go(fn func() any) *Thunk[any] {
	g.wg.Add(1)
	th := &Thunk[any]{}
	done := make(chan struct{})
	go func() {
		defer g.wg.Done()
		th.val = fn()
		close(done)
	}()
	th.fn = func() any {
		<-done
		return th.val
	}
	return th
}

// Wait blocks until all goroutines registered with this group complete.
func (g *BocGroup) Wait() { g.wg.Wait() }
