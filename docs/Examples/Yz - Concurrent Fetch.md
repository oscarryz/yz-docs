#example

## Concurrent HTTP Fetch

A classical concurrency benchmark: four HTTP requests that each take 2 seconds.

- **Sequential** execution would take ~8 seconds (requests one after the other).
- **Yz** takes ~2 seconds — all four fire immediately as concurrent bocs.

No explicit threads, no `Promise.all`, no `async/await`, no `goroutine` keywords.
Concurrency is automatic: every boc call returns a thunk that runs in the background.

```js
// Sequential execution would take ~8 seconds.
// Yz thunks fire all four immediately — total time ~2 seconds.

main: {
    print("start: ${time.now()}")

    // All four fire immediately concurrently.
    // a and b use inferred type; c and d declare the type explicitly.
    // Either way the type is String — the thunk is an invisible implementation detail.
    a: http.get("https://httpbin.org/delay/2")
    b: http.get("https://httpbin.org/delay/2")
    c String = http.get("https://httpbin.org/delay/2")
    d String = http.get("https://httpbin.org/delay/2")

    // Thunks are forced (materialized) here.
    // By the time we reach print(a), b, c, d are likely already done too.
    print("=== A ===")
    print(a)
    print("=== B ===")
    print(b)
    print("=== C ===")
    print(c)
    print("=== D ===")
    print(d)

    print("done:  ${time.now()}")
}
```

### How it works

Every assignment that calls a boc (including built-ins like `http.get`) immediately
spawns concurrently and returns a **thunk** — a lazy handle to the future result.

```
a: http.get("https://httpbin.org/delay/2")
```

This line does **not** wait for the response. It fires the request and moves on.
All four requests are in flight by the time the code reaches the first `print`.

**The type of `a` is still `String`**, not "thunk" type or "future" or "promise".
The thunk is an implementation detail invisible to the Yz programmer — you can
annotate the type explicitly and it remains the value type:

```
c String = http.get("https://httpbin.org/delay/2")
```

`c` is a `String`. It just happens to be computed concurrently.

A thunk is **forced** the moment the value is actually needed — here, when passed to `print`:

```
print(a)   // forces a — waits only if not yet done
```

Because all four were started at the same time, forcing them in sequence costs only
the time of the slowest one, not the sum of all four.


### String interpolation with inline calls

`time.now()` can be called directly inside a string interpolation:

```js
print("start: ${time.now()}")
```

The call is inlined — no intermediate variable needed. The thunk is forced
immediately because the string value is needed on that line.

### Output (example)

```
start: 2026-04-17T10:20:06+02:00
=== A ===
{ ... "url": "https://httpbin.org/delay/2" }
=== B ===
{ ... "url": "https://httpbin.org/delay/2" }
=== C ===
{ ... "url": "https://httpbin.org/delay/2" }
=== D ===
{ ... "url": "https://httpbin.org/delay/2" }
done:  2026-04-17T10:20:09+02:00
```

Three seconds elapsed for four two-second requests.
