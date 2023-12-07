```javascript
stack: {

	new: {
		Stack() 
	}
	Stack { 
		values [] // generic array
		push: {
			v  // generic value
			values << v
		}
		<< : { v push v }
		pop: {
			values.len() == 0 ? {
			   option.None()
			} { 
				option.Some(value.pop())
			}
		}
	}
}
main: {
	s: stack.new()
	s.push 42
	s << 43 
	// This is an error, as the stack inner `values` is now bound to `Int` types
	stack.push 'Oh no!'

    // Trying to use a generic type before is bound would return an error
    t: stack.new()
    x: t.pop() // not known yet, probably compile error
	10 + x.get() // yeah... probably not... or maybe, there's no way to know


    as_param: {
	    t stack.Stack 
		10 + t.pop().get() // would require t.value to be bound to Int // :? maybe
    }

}
```