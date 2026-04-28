#solved 

See:[Concurrency](https://github.com/oscarryz/yz-docs/blob/main/docs%2FFeatures%2FConcurrency.md)


----

How to make "transactions" without locking ala BOCs (Behavior Oriented stuff )

See: https://dl.acm.org/doi/pdf/10.1145/3622852


Implementing Behaviour-Oriented Concurrency (BoC) in Go is actually a great fit because Go’s runtime already handles the heavy lifting of green threads (goroutines). To adapt BoC to your "every function is an actor" paradigm, you need to solve one specific problem: **Atomic Multi-Resource Acquisition.**

In the BoC paper, a "Behaviour" won't run until it has locked all the "Cowns" it needs. In your compiler, this means your generated Go code shouldn't just spawn a goroutine; it should **enqueue a task** that only fires when its dependencies are ready.

Here is a conceptual architecture for your library/compiler backend.

---

## 1. The Core Abstraction: The `Cown`

A `Cown` (Concurrent Owner) in Go can be implemented as a struct that wraps a resource and manages a **request queue**.

Go

```
type Cown[T any] struct {
    data  T
    // A channel acting as a FIFO queue for behaviors wanting access
    requests chan chan bool 
}

func NewCown[T any](val T) *Cown[T] {
    c := &Cown[T]{
        data:     val,
        requests: make(chan chan bool, 100),
    }
    go c.serve()
    return c
}

func (c *Cown[T]) serve() {
    for req := range c.requests {
        req <- true  // Grant access
        <-req        // Wait for the behavior to finish (release)
    }
}
```

---

## 2. The `When` Scheduler (The "Brain")

To prevent deadlocks (a core BoC promise), the resources must be acquired in a **strict global order**. You can use the memory address of the Cowns as the sorting key.

Go

```
func When(cowns []Resource, behavior func()) {
    go func() {
        // 1. Sort cowns by address to prevent circular wait (Deadlock Prevention)
        sort.Slice(cowns, func(i, j int) bool {
            return reflect.ValueOf(cowns[i]).Pointer() < reflect.ValueOf(cowns[j]).Pointer()
        })

        // 2. Acquire all
        releases := make([]chan bool, len(cowns))
        for i, c := range cowns {
            releases[i] = c.acquire()
        }

        // 3. Execute the function
        behavior()

        // 4. Release all (in reverse order)
        for i := len(releases) - 1; i >= 0; i-- {
            releases[i] <- true
        }
    }()
}
```

---

## 3. Integrating with Your Compiler

If your language treats every function call as an actor, your compiler can analyze which variables are being passed. If a variable is "mutable" or "owned," you wrap it in a Cown.

### Your Language Code:

Pony

```
// Assume account1 and account2 are 'Cowns'
transfer(account1, account2, 100)
```

### Generated Go Code:

The backend doesn't just call a function; it wraps the call in a `When` block.

Go

```
// The compiler identifies that 'acc1' and 'acc2' are shared resources
lib.When([]lib.Resource{acc1, acc2}, func() {
    // This code only runs when the scheduler has locked BOTH accounts
    Transfer(acc1.Data, acc2.Data, 100)
})
```

---

## 4. Addressing Your "Everything is Concurrent" Goal

Since you want `foo()` and `bar()` to run concurrently by default, your Go backend should transform:

- **Function definitions** into "Behaviour" templates.
    
- **Function calls** into "Tasks" submitted to the BoC scheduler.
    

### The "Internal Complexity" Trade-off

To keep the user-facing side simple:

1. **Implicit Cowns:** Every object created in your language is automatically wrapped in a `Cown` by the Go backend.
    
2. **Static Analysis:** Your compiler determines the "Footprint" (which Cowns a function touches).
    
3. **No `go` keyword:** The user never sees `go`. Every function call is generated as a `When()` call. If a function touches zero resources, it's just a `go func()`. If it touches resources, the BoC scheduler handles the queue.
    

### Why this is better than "Standard" Actors:

In a standard Actor model, `transfer(a, b)` requires `a` to send a message to `b`, or a third coordinator. In your BoC-inspired Go backend, the `When` logic ensures that the `Transfer` logic sees both resources as local and safe, because the "Stop-the-world" lock for those specific two objects is handled by the library.

---

## Implementation Tips for Go

- **The Scheduler:** Use a worker pool if the number of actors becomes massive (millions), though goroutines are cheap enough that a `go func()` per behavior is usually fine for a prototype.
    
- **Memory Management:** Since Go has a GC, you don't need Pony's complex **ORCA** protocol for memory, but you _do_ need to ensure users don't manually leak pointers outside of the `Cown` wrapper.
    
	- **Channels as Queues:** Using a channel of channels (as shown in the `Cown` struct above) is the most idiomatic way to implement a "fair" FIFO acquisition queue in Go.