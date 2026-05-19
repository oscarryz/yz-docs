#example

[Notes on 1ML](https://shonfeder.github.io/themata/programming/notes-on-1ml.html)

A boc interface can be declared on its own — without a body — and the implementation assigned separately. This separates the contract from the concrete definition.

The signature declares what a `stack` must provide: an `empty` constructor and a `Stack` type with the expected methods. The implementation satisfies it structurally; no explicit declaration of conformance is needed. The interface can live in one file and the implementation in another.

```js
// Signature / interface
stack #(
	empty #(T, Stack(T))
	
	Stack #(
		T
		push     #(a T)
		pop      #(Option(T))
		is_empty #(Bool)
		size     #(Int)
	)
)

// usage
s: stack.empty(String)
s.push("hola")
s.is_empty() // false
s.size() // 1


// implementation
stack = {
	empty =  {
		T
		Stack(T)
	}
	Stack = {
		T
		array: [T]()
		push: {
			a T
			array.push(a)
		}
		pop: {
			array.length() > 0 ? {
				v: Option.Some(array[array.length()-1])
				array.remove(array.length() - 1)
				v
			}, {
				Option.None()
			}
		}

		is_empty: {
			array.length() == 0
		}
		size: {
			array.length()
		}
	}
}
```
