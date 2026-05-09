package rt

import "sync"

// Cown is the runtime representation of a concurrent owner.
// Every singleton boc struct embeds one. Method bodies run exclusively
// while holding the cown's mutex, preventing data races on struct fields.
type Cown struct{ mu sync.Mutex }

// Schedule runs fn while exclusively holding c, then releases c.
// Returns a Thunk that resolves once fn completes.
//
// fn must not block waiting for another behaviour that also needs c —
// use the split-BocGroup pattern (BocGroup declared outside Schedule,
// BocGroup.Wait() called after Schedule completes) to avoid deadlock.
func Schedule[T any](c *Cown, fn func() T) *Thunk[T] {
	return Go(func() T {
		c.mu.Lock()
		defer c.mu.Unlock()
		return fn()
	})
}
