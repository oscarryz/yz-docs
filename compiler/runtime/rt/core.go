package rt

import (
	"encoding/json"
	"fmt"
	"os"
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
// BocGroup — structured concurrency for a single boc invocation
// ---------------------------------------------------------------------------

// BocGroup collects thunks that must be forced before the enclosing boc
// returns. Underlying goroutines are already running (started by Schedule
// inside the cown), so sequential forcing loses no parallelism — total
// elapsed time equals the slowest goroutine.
type BocGroup struct {
	pending []func()
}

// Add registers fn to be called during Wait().
func (g *BocGroup) Add(fn func()) {
	g.pending = append(g.pending, fn)
}


// Wait forces all registered thunks sequentially.
func (g *BocGroup) Wait() {
	for _, fn := range g.pending {
		fn()
	}
}
