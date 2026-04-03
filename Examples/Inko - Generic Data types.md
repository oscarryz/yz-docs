#example

```js
stack: {

	new: {
		V
		Stack(V)
	}
	Stack: {
		V
		values [V]() // generic array
		push: {
			v V // generic value
			values << v
		}
		<<: { v V; push(v) }
		pop: {
			values.length() == 0 ? {
			   Option.None()
			}, {
				Option.Some(values.pop())
			}
		}
	}
}
main: {
	s: stack.new(Int)
	s.push(42)
	s << 43
	// This is an error, as the stack inner `values` is now bound to `Int` types
	stack.push('Oh no!')

    // Trying to use a generic type before is bound would an error
    t: stack.new(Int)
    x: t.pop() // runtime error, stack underflow.
	// x in an Option(Int)
	10 + x.get() // ok


    as_param: {
	    t stack.Stack(T)
		10 + t.pop().get() // would require t.value to be bound to Int // :? maybe
    }
}
```