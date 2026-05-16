// Trace harness for the closure_concurrent example.
//
// This is the hand-written Go equivalent of main.yz, instrumented with
// runtime/trace so that concurrent BOC execution can be visualised and
// compared across runtime/codegen changes.
//
// Run:
//   go run . && go tool trace trace.out
//
// What to look for in the viewer (Proc timeline):
//   - Three wide parallel bars for apply(b1=10), apply(b2=20), apply(b3=30)
//     — they hold distinct cowns so there is no contention.
//   - apply(b1=99) starts only after apply(b1=10) ends — b1.Cown is shared.
//   - Each "closure call" region is nested inside its "ScheduleMulti body"
//     region with no gap — the sync body (box.set) runs inline without
//     re-acquiring the cown.
//
// Future compile-time idea:
//   compile-time: [Tracer]
//   output: 'trace.out'
//
// When that infostring is supported by yzc, this harness can be retired and
// replaced by annotating main.yz directly.
package main

import (
	"context"
	"fmt"
	"os"
	"runtime/trace"
	"time"

	std "yz/runtime/rt"
)

// busyWork spins for ~5ms so the goroutine scheduler actually overlaps the
// three independent applies in the trace. Without it the work is too short
// and they appear sequential even though no cown contention exists.
func busyWork() {
	start := time.Now()
	for time.Since(start) < 5*time.Millisecond {
	}
}

// ── Box ──────────────────────────────────────────────────────────────────────
// Corresponds to: Box: { val Int; set #(v Int) { val = v } }

type Box struct {
	std.Cown
	val std.Int
}

func NewBox(val std.Int) *Box { return &Box{val: val} }

func (self *Box) set(v std.Int) std.Unit {
	self.val = v
	return std.TheUnit
}

func (self *Box) Set(v std.Int) std.Unit {
	return std.LazyUnit(std.Schedule(&self.Cown, func() std.Unit {
		return self.set(v)
	}))
}

// ── apply ─────────────────────────────────────────────────────────────────────
// Corresponds to: apply #(a Box, fn #()) { fn() }
//
// The lowerer detects that fn closes over a Box param (a), so calls to
// a.set() inside fn are emitted as the sync body (lowercase) rather than
// the async uppercase+Force path. Without that, fn() would deadlock:
// ScheduleMulti already holds a.Cown and Set() would try to acquire it again.

type _applyBoc struct {
	std.Cown
	a  *Box
	fn func() std.Unit
}

func (self *_applyBoc) call(label string, a *Box, fn func() std.Unit) std.Unit {
	taskCtx, task := trace.NewTask(context.Background(), "apply("+label+")")
	_bg0 := &std.BocGroup{}
	_sched := std.ScheduleMulti([]*std.Cown{&self.Cown, &a.Cown}, func() std.Unit {
		trace.WithRegion(taskCtx, "ScheduleMulti body", func() {
			self.a = a
			self.fn = fn
			trace.WithRegion(taskCtx, "closure call", func() {
				self.fn() // sync — no cown re-acquisition
			})
		})
		return std.TheUnit
	})
	return std.LazyUnit(std.NewThunk(func() std.Unit {
		_sched.Force()
		_bg0.Wait()
		task.End()
		return std.TheUnit
	}))
}

// ── main ──────────────────────────────────────────────────────────────────────

type _mainBoc struct{ std.Cown }

func (self *_mainBoc) Call() std.Unit {
	return std.LazyUnit(std.NewThunk(func() std.Unit {
		_bg0 := &std.BocGroup{}
		var b1, b2, b3 *Box

		std.Schedule(&self.Cown, func() std.Unit {
			b1 = NewBox(std.NewInt(0))
			b2 = NewBox(std.NewInt(0))
			b3 = NewBox(std.NewInt(0))

			// Three independent applies — parallel, no shared cowns.
			_st0 := (&_applyBoc{}).call("b1=10", b1, func() std.Unit {
				busyWork()
				return b1.set(std.NewInt(10))
			})
			_bg0.GoWait(_st0)

			_st1 := (&_applyBoc{}).call("b2=20", b2, func() std.Unit {
				busyWork()
				return b2.set(std.NewInt(20))
			})
			_bg0.GoWait(_st1)

			_st2 := (&_applyBoc{}).call("b3=30", b3, func() std.Unit {
				busyWork()
				return b3.set(std.NewInt(30))
			})
			_bg0.GoWait(_st2)

			// Shared cown: also targets b1 — serialized after apply(b1=10).
			_st3 := (&_applyBoc{}).call("b1=99", b1, func() std.Unit {
				busyWork()
				return b1.set(std.NewInt(99))
			})
			_bg0.GoWait(_st3)

			return std.TheUnit
		}).Force()

		_bg0.Wait()

		fmt.Printf("b1.val = %v\n", b1.val) // 99
		fmt.Printf("b2.val = %v\n", b2.val) // 20
		fmt.Printf("b3.val = %v\n", b3.val) // 30
		return std.TheUnit
	}))
}

var Main = &_mainBoc{}

func main() {
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := trace.Start(f); err != nil {
		panic(err)
	}
	defer trace.Stop()

	Main.Call().Force()
}
