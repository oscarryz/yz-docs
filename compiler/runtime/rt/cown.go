package rt

import (
	"runtime"
	"sync/atomic"
)

// request is a node in a cown's pending-behaviour queue.
type request struct {
	next atomic.Pointer[request]
	beh  *behaviour
}

// behaviour is a unit of work scheduled on one or more cowns.
// It runs when every required cown has granted it a token (count reaches 0).
type behaviour struct {
	count atomic.Int64
	run   func()
}

// Cown is the runtime representation of a concurrent owner.
// Every singleton boc struct embeds one. A nil `last` means the cown is idle.
type Cown struct {
	last atomic.Pointer[request]
}

// Schedule runs fn while exclusively holding c, then releases c.
// Returns a Thunk that resolves once fn completes.
//
// fn must not block waiting for another behaviour that also needs c —
// use the split-BocGroup pattern (BocGroup declared outside Schedule,
// BocGroup.Wait() called after Schedule completes) to avoid deadlock.
func Schedule[T any](c *Cown, fn func() T) *Thunk[T] {
	done := make(chan struct{})
	var result T

	b := &behaviour{}
	b.count.Store(1)
	req := &request{beh: b}
	b.run = func() {
		result = fn()
		close(done)
		releaseCown(c, req)
	}

	enqueueCown(c, req)

	return NewThunk(func() T {
		<-done
		return result
	})
}

// enqueueCown atomically adds req to c's pending-behaviour queue.
// If c was idle, grants req's behaviour a token immediately.
func enqueueCown(c *Cown, req *request) {
	prev := c.last.Swap(req)
	if prev == nil {
		// Cown was idle — grant token; run if all cowns have granted.
		if req.beh.count.Add(-1) == 0 {
			go req.beh.run()
		}
	} else {
		// Cown busy — link req as prev's successor.
		prev.next.Store(req)
	}
}

// releaseCown is called after a behaviour's fn completes.
// It either marks c idle or passes the token to the next waiting behaviour.
func releaseCown(c *Cown, req *request) {
	// Try to mark cown idle: we're still the tail.
	if c.last.CompareAndSwap(req, nil) {
		return
	}
	// A new request was enqueued after us — wait for it to link itself.
	var next *request
	for {
		if next = req.next.Load(); next != nil {
			break
		}
		runtime.Gosched()
	}
	// Grant token to successor; launch it if all its cowns have now granted.
	if next.beh.count.Add(-1) == 0 {
		go next.beh.run()
	}
}
